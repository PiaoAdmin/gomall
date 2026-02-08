# 自动下单Agent

基于 LangGraph 实现的智能购物助手，支持完整的电商下单流程。

## 功能特性

- 🔍 **智能商品搜索** - 自然语言搜索商品
- 🛒 **购物车管理** - 添加/移除商品、查看购物车
- 📦 **一键下单** - 自动收集收货信息并完成下单
- 🤖 **多轮交互** - 使用LangGraph interrupt机制实现人机协作

## 工作流程

```
用户输入需求
    ↓
搜索商品 (search)
    ↓
🛑 确认选择 (confirm_sku) - interrupt点
    ↓
处理选择 (handle_sku)
    ├─ 继续购物 → 搜索商品
    ├─ 查看购物车 → view_cart
    └─ 去结算 → view_cart
         ↓
查看购物车 (view_cart)
    ↓
🛑 确认购物车 (confirm_cart) - interrupt点
    ↓
处理确认 (handle_cart)
    ├─ 继续购物 → 搜索商品
    └─ 去结算 → collect_address
         ↓
🛑 收集地址 (collect_address) - interrupt点
    ↓
处理地址 (handle_address)
    ├─ 信息不全 → 继续收集
    └─ 信息完整 → confirm_order
         ↓
🛑 确认订单 (confirm_order) - interrupt点
    ↓
处理确认 (handle_order)
    ├─ 修改地址 → collect_address
    └─ 确认下单 → place_order
         ↓
执行下单 (place_order)
    ↓
结束 (END)
```

## 快速开始

### 1. 环境准备

```bash
# 安装依赖（如需要）
cd /home/piao/workspace/piao_blog/p_mall/agent
source .venv/bin/activate  # 使用已有虚拟环境
pip install python-dotenv  # 如果未安装

# 配置环境变量
cp .env.example .env
# 编辑.env文件，填入你的API Key
```

**.env 文件示例：**
```bash
OPENAI_API_KEY=sk-your-api-key-here
OPENAI_BASE_URL=https://dashscope.aliyuncs.com/compatible-mode/v1
OPENAI_MODEL=qwen-plus
```

### 2. 运行Agent

```bash
cd /home/piao/workspace/piao_blog/p_mall/agent
source .venv/bin/activate
python -m auto_order_agent.main
```

### 3. 使用示例

```
👤 你: 我想买红米手机

🤖 Agent: 🛍️ 搜索结果：

【1】红米K30 256GB+12GB版
   价格: ¥1999.99（原价 ¥2399.99）
   库存: 100件
   SKU ID: 123
   
【2】红米K60 512GB+16GB版
   价格: ¥2499.99
   库存: 50件
   SKU ID: 124

请选择商品编号

👤 你: 1

🤖 Agent: 已添加到购物车: 红米K30 256GB+12GB版 x 1
商品已加入购物车，是否继续购物？
- 输入「查看购物车」查看当前购物车
- 输入「去结算」开始下单流程

👤 你: 去结算

🤖 Agent: 🛒 购物车详情：

1. 红米K30 256GB+12GB版
   数量: 1
   单价: ¥1999.99
   
总数量: 1件
总金额: ¥1999.99

确认购物车内容？输入「去结算」继续

👤 你: 去结算

🤖 Agent: 好的，开始收集收货信息。

请提供以下信息：
1. 收货人姓名
2. 详细地址
3. 城市
4. 邮编
5. 联系邮箱

👤 你: 姓名：张三
地址：北京市朝阳区xx街xx号
城市：北京
邮编：100000
邮箱：user@example.com

🤖 Agent: 📦 收货信息已收集：

收货人: 张三
地址: 北京市朝阳区xx街xx号
城市: 北京
邮编: 100000
邮箱: user@example.com

请确认收货信息是否正确？（输入「确认」继续）

👤 你: 确认

🤖 Agent: ✅ 订单创建成功！订单号: ORDER_123456789

会话结束。继续购物或输入 'exit' 退出。
```

## 技术架构

### API封装 (api_client.py)

```python
class PmallOrderAPIClient:
    - login()                # 用户登录
    - search_products()      # 搜索商品
    - get_product_detail()   # 商品详情
    - get_cart()             # 获取购物车
    - add_to_cart()          # 添加到购物车
    - remove_from_cart()     # 移除商品
    - place_order()          # 下单
    - list_orders()          # 订单列表
```

### LangChain工具 (tools.py)

6个工具供LLM调用：
- `search_products_tool` - 搜索商品
- `get_product_detail_tool` - 查看详情
- `view_cart_tool` - 查看购物车
- `add_to_cart_tool` - 添加到购物车
- `remove_from_cart_tool` - 移除商品
- `place_order_tool` - 执行下单

### 状态定义 (state.py)

```python
class AutoOrderState(TypedDict):
    messages: list[BaseMessage]            # 消息历史
    search_results: Optional[List[Dict]]   # 搜索结果
    selected_sku: Optional[Dict]           # 选中的SKU
    cart_items: Optional[List[Dict]]       # 购物车
    shipping_address: Optional[Dict]       # 收货地址
    order_id: Optional[str]                # 订单号
    next_step: Literal[...]                # 下一步
```

### LangGraph工作流 (graph.py)

**11个节点**：
1. `search` - 搜索商品
2. `confirm_sku` - 确认SKU（interrupt）
3. `handle_sku` - 处理SKU选择
4. `view_cart` - 查看购物车
5. `confirm_cart` - 确认购物车（interrupt）
6. `handle_cart` - 处理购物车确认
7. `collect_address` - 收集地址（interrupt）
8. `handle_address` - 处理地址
9. `confirm_order` - 确认订单（interrupt）
10. `handle_order` - 处理订单确认
11. `place_order` - 执行下单

**4个interrupt点**：
- `confirm_sku` - 商品选择后暂停
- `confirm_cart` - 购物车确认前暂停
- `collect_address` - 地址收集时暂停
- `confirm_order` - 订单确认前暂停

## 核心特性

### 1. Interrupt机制实现多轮交互

```python
# 在关键节点前自动暂停
graph = workflow.compile(
    checkpointer=MemorySaver(),
    interrupt_before=["confirm_sku", "confirm_cart", "collect_address", "confirm_order"]
)

# 检测interrupt状态
snapshot = graph.get_state(config)
if snapshot.next and "confirm_sku" in snapshot.next:
    # 获取用户输入并注入
    graph.update_state(config, {"messages": [...]}, as_node="confirm_sku")
    # 继续执行
    graph.stream(None, config)
```

### 2. 条件路由

每个处理节点根据用户意图返回不同的`next_step`，通过条件路由跳转：

```python
def route_after_sku_handling(state):
    if state["next_step"] == "view_cart":
        return "view_cart"
    elif state["next_step"] == "search":
        return "search"
    else:
        return "confirm_sku"
```

### 3. LLM工具调用

```python
llm_with_tools = llm.bind_tools([
    search_products_tool,
    add_to_cart_tool,
    view_cart_tool,
    ...
])

response = llm_with_tools.invoke(messages)
# LLM自动选择合适的工具并调用
```

## 与product_listing_agent的对比

| 特性 | product_listing | auto_order |
|-----|----------------|-----------|
| 节点数 | 5 | 11 |
| interrupt点 | 1 (confirm) | 4 (多个确认点) |
| 工具数 | 3 | 6 |
| 交互轮次 | 1-2轮 | 4-5轮 |
| 流程复杂度 | 简单（单一流程） | 复杂（多分支） |
| 状态字段 | 6个 | 7个 |

## 开发说明

### 添加新功能

1. **新增API接口** - 在`api_client.py`中添加方法
2. **创建工具** - 在`tools.py`中用`@tool`装饰器封装
3. **添加节点** - 在`graph.py`中定义节点函数
4. **配置路由** - 使用`add_conditional_edges`设置流转
5. **设置interrupt** - 在`interrupt_before`列表中添加节点名

### 调试技巧

```python
# 打印当前状态
snapshot = graph.get_state(config)
print(f"Current nodes: {snapshot.next}")
print(f"State: {snapshot.values}")

# 查看消息历史
for msg in state["messages"]:
    print(f"{msg.__class__.__name__}: {msg.content}")
```

## 常见问题

**Q: 为什么需要这么多interrupt点？**

A: 每个interrupt点代表一次人机交互确认。购物流程需要多次确认：选商品、确认购物车、填地址、确认订单，确保用户掌控全流程。

**Q: 如何处理用户中途改变主意？**

A: 每个`handle_*`节点会检测用户意图，支持跳回上一步或切换流程。例如在确认订单时说"修改地址"会跳回`collect_address`。

**Q: LLM会自动调用工具吗？**

A: 是的。使用`llm.bind_tools()`后，LLM会根据上下文自动选择合适的工具并调用。

## License

MIT
