# LangGraph 商品上架Agent架构说明

## 核心特性

基于 **LangGraph 0.2** 实现的人机协作工作流，使用 **interrupt_before + MemorySaver** 模式实现多轮交互。

## 状态定义

```python
from typing import TypedDict, Literal, Optional
from langchain_core.messages import BaseMessage

class ProductListingState(TypedDict):
    """状态类型继承TypedDict，LangGraph自动合并更新"""
    messages: list[BaseMessage]          # 消息历史
    product_data: Optional[dict]         # 商品数据
    validation_status: str               # 验证状态
    retry_count: int                     # 重试次数
    error_message: str                   # 错误信息
    next_step: Literal["complete_info", "confirm", "validate", "create", "retry", "end"]
```

**关键点**：
- 继承 `TypedDict` 而非普通 dict
- LangGraph 自动合并节点返回的部分状态（类似 React setState）
- `next_step` 用于条件路由决策

## 工作流节点

```python
from langgraph.graph import StateGraph, END

workflow = StateGraph(ProductListingState)

# 5个核心节点
workflow.add_node("complete_info", complete_product_info_node)  # LLM补全商品信息
workflow.add_node("confirm", user_confirmation_node)            # 用户确认（interrupt点）
workflow.add_node("validate", validate_user_input_node)         # 验证用户响应
workflow.add_node("create", create_product_node)                # 调用API创建
workflow.add_node("retry", retry_with_fix_node)                 # 失败重试
```

## 人机交互机制

### 1. Interrupt Before 模式

```python
from langgraph.checkpoint.memory import MemorySaver

memory = MemorySaver()  # 状态持久化
graph = workflow.compile(
    checkpointer=memory,
    interrupt_before=["confirm"]  # 在confirm节点前暂停
)
```

**原理**：
- Graph 在执行到 `confirm` 节点**之前**自动暂停
- 状态保存到 checkpointer（内存中）
- 等待外部输入

### 2. 状态注入与恢复

```python
# 初次运行：complete_info节点生成商品数据后暂停
config = {"configurable": {"thread_id": "1"}}
for event in graph.stream(initial_state, config):
    # 打印商品信息，询问用户确认
    ...

# 检测是否在interrupt点
snapshot = graph.get_state(config)
if snapshot.next and "confirm" in snapshot.next:
    # 获取用户输入
    user_input = input("👤 你: ")
    
    # 注入用户消息到confirm节点
    graph.update_state(
        config,
        {"messages": [HumanMessage(content=user_input)]},
        as_node="confirm"  # 以confirm节点身份更新状态
    )
    
    # 从checkpoint恢复执行（传入None表示继续）
    for event in graph.stream(None, config):
        ...
```

**关键API**：
- `graph.get_state(config)` - 获取当前状态快照，检查 `snapshot.next`
- `graph.update_state(config, values, as_node)` - 注入新数据
- `graph.stream(None, config)` - 从checkpoint继续执行

### 3. 条件路由

```python
def route_after_validation(state: ProductListingState) -> str:
    """验证后路由：确认创建 or 继续修改"""
    if state["next_step"] == "create":
        return "create"
    return "confirm"  # 循环回confirm，再次interrupt

workflow.add_conditional_edges(
    "validate",
    route_after_validation,
    {"create": "create", "confirm": "confirm"}
)
```

## 完整流程图

```
用户输入 "红米K30"
    ↓
complete_info (LLM生成完整商品数据)
    ↓
🛑 interrupt_before=["confirm"] (暂停)
    ↓
用户确认/修改 → update_state注入消息
    ↓
confirm (透传节点: return state)
    ↓
validate (判断用户意图)
    ↓
  ┌─────┴─────┐
  ↓           ↓
create     confirm (修改循环)
  ↓           ↓
retry      🛑 interrupt (再次暂停)
  ↓
 END
```

## 核心节点实现

### complete_info: LLM补全

```python
def complete_product_info_node(state: ProductListingState):
    llm = ChatOpenAI()
    llm_with_tools = llm.bind_tools([
        get_categories_tool,
        get_brands_tool
    ])
    
    response = llm_with_tools.invoke(state["messages"])
    
    # 解析JSON，格式化显示
    product_data = extract_json(response.content)
    display_text = format_product_data(product_data)
    
    return {
        "messages": [AIMessage(content=display_text + "\n请确认...")],
        "product_data": product_data,
        "next_step": "confirm"
    }
```

### confirm: 透传节点

```python
def user_confirmation_node(state: ProductListingState):
    """仅用作interrupt点，不做处理"""
    return state
```

### validate: 意图识别

```python
def validate_user_input_node(state: ProductListingState):
    last_message = state["messages"][-1].content
    
    if any(kw in last_message for kw in ["确认", "是", "好"]):
        return {"next_step": "create"}
    
    # 用户要修改，调用LLM更新数据
    updated_data = llm_update(state["product_data"], last_message)
    return {
        "product_data": updated_data,
        "messages": [AIMessage(content="已更新：" + format(updated_data))],
        "next_step": "confirm"  # 循环回去
    }
```

## LangGraph关键特性总结

| 特性 | 用法 | 作用 |
|-----|------|------|
| `StateGraph` | 定义状态机 | 类型安全的状态管理 |
| `TypedDict` | 状态类型 | 自动合并更新 |
| `MemorySaver` | checkpointer | 状态持久化 |
| `interrupt_before` | 暂停点 | 人机交互 |
| `update_state` | 注入数据 | 外部输入 |
| `stream(None, config)` | 恢复执行 | 从checkpoint继续 |
| `conditional_edges` | 条件路由 | 动态流程控制 |

## 为什么这样设计

1. **TypedDict状态** - 自动合并，节点只需返回变更字段
2. **interrupt_before** - 声明式暂停，无需手动判断
3. **MemorySaver** - 自动保存/恢复，支持多轮对话
4. **update_state + as_node** - 精确控制数据注入位置
5. **stream(None)** - 优雅恢复，避免重复执行

这套模式完全遵循LangGraph最佳实践，消除了手动状态管理和循环输出问题。
