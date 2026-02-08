"""LangGraph workflow for Auto Order Agent

实现完整的自动下单流程状态机
"""

import os
import json
from typing import Dict, Any
from langchain_openai import ChatOpenAI
from langchain_core.messages import HumanMessage, AIMessage, SystemMessage
from langgraph.graph import StateGraph, END
from langgraph.checkpoint.memory import MemorySaver

from .state import (
    AutoOrderState,
    SEARCH_SYSTEM_PROMPT,
    CONFIRM_SKU_SYSTEM_PROMPT,
    VIEW_CART_SYSTEM_PROMPT,
    COLLECT_ADDRESS_PROMPT,
    CONFIRM_ORDER_PROMPT
)
from .tools import get_all_tools


def get_llm():
    """获取LLM实例"""
    return ChatOpenAI(
        model=os.getenv("OPENAI_MODEL", "gpt-4"),
        temperature=0.7
    )


# ==================== 节点函数 ====================

def search_products_node(state: AutoOrderState) -> Dict[str, Any]:
    """节点1: 搜索商品并展示结果"""
    llm = get_llm()
    llm_with_tools = llm.bind_tools(get_all_tools())
    
    # 构建消息
    messages = [SystemMessage(content=SEARCH_SYSTEM_PROMPT)] + state["messages"]
    
    # 调用LLM
    response = llm_with_tools.invoke(messages)
    
    # DEBUG: 打印LLM响应
    print(f"\n🔍 DEBUG - LLM响应类型: {type(response)}")
    print(f"🔍 DEBUG - 有tool_calls属性: {hasattr(response, 'tool_calls')}")
    if hasattr(response, 'tool_calls'):
        print(f"🔍 DEBUG - tool_calls内容: {response.tool_calls}")
    print(f"🔍 DEBUG - 响应内容: {response.content[:200] if response.content else 'None'}...\n")
    
    # 处理工具调用
    new_messages = state["messages"].copy()
    search_results = []
    
    if hasattr(response, 'tool_calls') and response.tool_calls:
        # 手动执行工具
        print("✅ 检测到工具调用，开始执行...")
        tools_map = {tool.name: tool for tool in get_all_tools()}
        results = []
        
        for tool_call in response.tool_calls:
            tool_name = tool_call["name"]
            print(f"  📞 调用工具: {tool_name}, 参数: {tool_call['args']}")
            if tool_name in tools_map:
                try:
                    result = tools_map[tool_name].invoke(tool_call["args"])
                    print(f"  ✅ 工具执行成功，结果: {str(result)[:100]}...")
                    results.append(f"{tool_name}结果:\n{result}")
                    
                    # 保存搜索结果
                    if tool_name == "search_products_tool":
                        import json
                        try:
                            search_results = json.loads(result)
                        except:
                            pass
                except Exception as e:
                    print(f"  ❌ 工具执行失败: {str(e)}")
                    results.append(f"{tool_name}执行失败: {str(e)}")
        
        # 将工具结果组合成一条消息，让LLM生成用户友好的回复
        tool_results = "\n\n".join(results)
        final_messages = messages + [
            AIMessage(content=f"工具执行完成，结果如下:\n{tool_results}\n\n请基于以上结果，以友好的方式展示商品信息给用户。")
        ]
        
        final_response = llm.invoke(final_messages)
        new_messages.append(final_response)
    else:
        # 没有工具调用，直接使用LLM回复
        print("⚠️ 未检测到工具调用，使用LLM直接回复")
        new_messages.append(response)
    
    return {
        "messages": new_messages,
        "search_results": search_results if search_results else state.get("search_results"),
        "next_step": "confirm_sku"
    }


def confirm_sku_selection_node(state: AutoOrderState) -> Dict[str, Any]:
    """节点2: 确认SKU选择（interrupt点）"""
    return state


def handle_sku_selection_node(state: AutoOrderState) -> Dict[str, Any]:
    """节点3: 处理SKU选择 - 解析用户输入的序号并添加到购物车"""
    last_message = state["messages"][-1].content.strip()
    
    # 检查用户意图
    if any(keyword in last_message.lower() for keyword in ["查看购物车", "购物车", "查看", "cart"]):
        return {"next_step": "view_cart"}
    
    if any(keyword in last_message.lower() for keyword in ["去结算", "结算", "checkout", "下单"]):
        return {"next_step": "view_cart"}  # 先查看购物车
    
    if any(keyword in last_message.lower() for keyword in ["继续购物", "再看看"]):
        return {"next_step": "search"}
    
    # 检查是否为数字选择（如 "1", "2", "3"）
    if last_message.isdigit():
        choice_idx = int(last_message) - 1  # 用户说1表示第0个
        search_results = state.get("search_results", [])
        
        if search_results and 0 <= choice_idx < len(search_results):
            selected_sku = search_results[choice_idx]
            sku_id = selected_sku.get("sku_id")
            sku_name = selected_sku.get("sku_name", "未知商品")
            
            print(f"📋 用户选择第{last_message}个商品: {sku_name} (SKU ID: {sku_id})")
            
            # 直接调用API添加到购物车
            from .tools import add_to_cart_tool
            result = add_to_cart_tool.invoke({"sku_id": sku_id, "quantity": 1})
            
            new_messages = state["messages"].copy()
            new_messages.append(AIMessage(content=result))
            
            return {
                "messages": new_messages,
                "selected_sku": selected_sku,
                "next_step": "view_cart"
            }
        else:
            new_messages = state["messages"].copy()
            new_messages.append(AIMessage(content=f"❌ 无效选择，请输入1-{len(search_results)}之间的数字"))
            return {"messages": new_messages, "next_step": "confirm_sku"}
    
    # 否则调用LLM处理其他情况
    llm = get_llm()
    llm_with_tools = llm.bind_tools(get_all_tools())
    
    messages = [SystemMessage(content=CONFIRM_SKU_SYSTEM_PROMPT)] + state["messages"]
    response = llm_with_tools.invoke(messages)
    
    new_messages = state["messages"].copy()
    
    if hasattr(response, 'tool_calls') and response.tool_calls:
        # 手动执行工具
        tools_map = {tool.name: tool for tool in get_all_tools()}
        results = []
        
        for tool_call in response.tool_calls:
            tool_name = tool_call["name"]
            if tool_name in tools_map:
                try:
                    result = tools_map[tool_name].invoke(tool_call["args"])
                    results.append(f"{result}")
                except Exception as e:
                    results.append(f"执行失败: {str(e)}")
        
        # 直接展示工具结果
        tool_results = "\n".join(results)
        new_messages.append(AIMessage(content=tool_results))
    else:
        new_messages.append(response)
    
    return {
        "messages": new_messages,
        "next_step": "confirm_sku"
    }


def view_cart_node(state: AutoOrderState) -> Dict[str, Any]:
    """节点4: 查看购物车"""
    llm = get_llm()
    llm_with_tools = llm.bind_tools(get_all_tools())
    
    messages = [SystemMessage(content=VIEW_CART_SYSTEM_PROMPT)] + state["messages"]
    response = llm_with_tools.invoke(messages)
    
    new_messages = state["messages"].copy()
    
    if hasattr(response, 'tool_calls') and response.tool_calls:
        tools_map = {tool.name: tool for tool in get_all_tools()}
        results = []
        
        for tool_call in response.tool_calls:
            tool_name = tool_call["name"]
            if tool_name in tools_map:
                try:
                    result = tools_map[tool_name].invoke(tool_call["args"])
                    results.append(f"{result}")
                except Exception as e:
                    results.append(f"执行失败: {str(e)}")
        
        tool_results = "\n\n".join(results)
        # 让LLM格式化展示
        final_messages = messages + [AIMessage(content=f"购物车数据:\n{tool_results}\n\n请格式化展示购物车内容")]
        final_response = llm.invoke(final_messages)
        new_messages.append(final_response)
    else:
        new_messages.append(response)
    
    return {
        "messages": new_messages,
        "next_step": "confirm_cart"
    }


def confirm_cart_node(state: AutoOrderState) -> Dict[str, Any]:
    """节点5: 确认购物车（interrupt点）"""
    return state


def handle_cart_confirmation_node(state: AutoOrderState) -> Dict[str, Any]:
    """节点6: 处理购物车确认"""
    last_message = state["messages"][-1].content.lower()
    
    if any(keyword in last_message for keyword in ["继续购物", "再逛逛", "添加"]):
        return {"next_step": "search"}
    
    if any(keyword in last_message for keyword in ["去结算", "结算", "确认", "下单"]):
        return {
            "messages": state["messages"] + [
                AIMessage(content="好的，开始收集收货信息。\n\n请提供以下信息：\n1. 收货人姓名\n2. 详细地址\n3. 城市\n4. 邮编\n5. 联系邮箱")
            ],
            "next_step": "collect_address"
        }
    
    # 其他输入，继续等待
    return {
        "messages": state["messages"] + [
            AIMessage(content='请输入「去结算」继续，或「继续购物」添加更多商品')
        ],
        "next_step": "confirm_cart"
    }


def collect_address_node(state: AutoOrderState) -> Dict[str, Any]:
    """节点7: 收集收货地址（interrupt点）"""
    return state


def handle_address_node(state: AutoOrderState) -> Dict[str, Any]:
    """节点8: 处理收货地址 - 使用LLM智能解析"""
    llm = get_llm()
    
    user_input = state["messages"][-1].content
    
    # 使用LLM智能提取地址信息
    extract_prompt = f"""从用户输入中提取收货地址信息。用户可能用各种格式提供信息。

用户输入: {user_input}

需要提取的字段：
1. name (收货人姓名) - 通常是人名，如"张三"、"李四"、"tbb"
2. street_address (详细地址) - 街道、小区、门牌号等，如"海淀区西土城路10号"、"朝阳区xx街"
3. city (城市) - 如"北京市"、"上海"、"深圳"
4. zip_code (邮编) - 6位数字，如100876、100012
5. email (邮箱) - 如test@qq.com、abc@example.com

常见格式示例：
- "张三 北京市 海淀区西土城路10号 100876 test@qq.com"
- "tb 北京市 海淀区 100012 1919456770@123.com 2222@123.com" (取第一个邮箱)
- "名字 tbb street address 西土城路" (tbb是姓名，西土城路是street_address)

规则：
- 提取所有能识别的字段
- 无法确定的字段设为null
- zip_code必须是整数
- 如果有多个邮箱，取第一个
- 返回纯JSON，格式: {{"name": "...", "street_address": "...", "city": "...", "zip_code": 数字, "email": "..."}}
- 不要包含markdown代码块标记```"""
    
    response = llm.invoke([HumanMessage(content=extract_prompt)])
    
    try:
        # 提取JSON
        content = response.content.strip()
        # 移除可能的markdown代码块
        if "```json" in content:
            start = content.find("```json") + 7
            end = content.find("```", start)
            json_str = content[start:end].strip()
        elif "```" in content:
            start = content.find("```") + 3
            end = content.find("```", start)
            json_str = content[start:end].strip()
        else:
            json_str = content
        
        import json
        address_data = json.loads(json_str)
        
        print(f"📍 解析地址结果: {address_data}")
        
        # 合并到现有地址（保留之前输入的字段）
        current_address = state.get("shipping_address", {})
        for key in ["name", "street_address", "city", "zip_code", "email"]:
            if key in address_data and address_data[key] is not None:
                current_address[key] = address_data[key]
        
        # 检查是否有缺失字段
        required_fields = ["name", "street_address", "city", "zip_code", "email"]
        missing_fields = [f for f in required_fields if not current_address.get(f)]
        
        if missing_fields:
            # 友好提示缺失字段
            field_names = {"name": "收货人姓名", "street_address": "详细地址", "city": "城市", "zip_code": "邮编", "email": "邮箱"}
            missing_msg = "还需要以下信息：\n" + "\n".join(f"- {field_names.get(f, f)}" for f in missing_fields)
            missing_msg += "\n\n请补充提供（可以一次性提供所有信息）"
            
            return {
                "messages": state["messages"] + [AIMessage(content=missing_msg)],
                "shipping_address": current_address,
                "next_step": "collect_address"
            }
        
        # 信息完整，展示并请求确认
        confirm_msg = f"""✅ 收货信息已收集完成：

📝 收货人: {current_address['name']}
📍 详细地址: {current_address['street_address']}
🏙️ 城市: {current_address['city']}
📮 邮编: {current_address['zip_code']}
📧 邮箱: {current_address['email']}

请确认收货信息是否正确？
- 输入「确认」或「对的」继续下单
- 或直接输入需要修改的信息"""
        
        return {
            "messages": state["messages"] + [AIMessage(content=confirm_msg)],
            "shipping_address": current_address,
            "next_step": "confirm_order"
        }
        
    except Exception as e:
        print(f"❌ 地址解析错误: {e}")
        return {
            "messages": state["messages"] + [
                AIMessage(content=f"抱歉，未能识别您的地址信息。\n\n请按以下格式提供：\n姓名 城市 详细地址 邮编 邮箱\n\n示例: 张三 北京市 海淀区中关村大街1号 100080 test@qq.com")
            ],
            "next_step": "collect_address"
        }


def confirm_order_node(state: AutoOrderState) -> Dict[str, Any]:
    """节点9: 确认订单（interrupt点）"""
    return state


def handle_order_confirmation_node(state: AutoOrderState) -> Dict[str, Any]:
    """节点10: 处理订单确认"""
    last_message = state["messages"][-1].content.lower()
    
    if any(keyword in last_message for keyword in ["确认", "下单", "ok", "好"]):
        return {"next_step": "place_order"}
    
    if "地址" in last_message or "修改" in last_message:
        return {
            "messages": state["messages"] + [AIMessage(content="请重新输入收货信息")],
            "shipping_address": None,
            "next_step": "collect_address"
        }
    
    return {
        "messages": state["messages"] + [AIMessage(content='请输入「确认」下单，或「修改」重新填写地址')],
        "next_step": "confirm_order"
    }


def place_order_node(state: AutoOrderState) -> Dict[str, Any]:
    """节点11: 执行下单"""
    from .tools import place_order_tool
    
    address = state["shipping_address"]
    
    try:
        result = place_order_tool.invoke({
            "email": address["email"],
            "name": address["name"],
            "street_address": address["street_address"],
            "city": address["city"],
            "zip_code": address["zip_code"]
        })
        
        return {
            "messages": state["messages"] + [AIMessage(content=f"✅ {result}")],
            "next_step": "end"
        }
    except Exception as e:
        return {
            "messages": state["messages"] + [AIMessage(content=f"❌ 下单失败: {str(e)}")],
            "next_step": "end"
        }


# ==================== 路由函数 ====================

def route_after_sku_handling(state: AutoOrderState) -> str:
    """SKU处理后路由"""
    next_step = state["next_step"]
    if next_step == "view_cart":
        return "view_cart"
    elif next_step == "search":
        return "search"
    else:
        return "confirm_sku"


def route_after_cart_handling(state: AutoOrderState) -> str:
    """购物车处理后路由"""
    if state["next_step"] == "search":
        return "search"
    elif state["next_step"] == "collect_address":
        return "collect_address"
    else:
        return "confirm_cart"


def route_after_address_handling(state: AutoOrderState) -> str:
    """地址处理后路由"""
    if state["next_step"] == "confirm_order":
        return "confirm_order"
    else:
        return "collect_address"


def route_after_order_confirmation(state: AutoOrderState) -> str:
    """订单确认后路由"""
    next_step = state["next_step"]
    if next_step == "place_order":
        return "place_order"
    elif next_step == "collect_address":
        return "collect_address"
    else:
        return "confirm_order"


def route_after_place_order(state: AutoOrderState) -> str:
    """下单后路由"""
    return "end"


# ==================== 构建工作流 ====================

def build_graph():
    """构建LangGraph工作流"""
    workflow = StateGraph(AutoOrderState)
    
    # 添加所有节点
    workflow.add_node("search", search_products_node)
    workflow.add_node("confirm_sku", confirm_sku_selection_node)
    workflow.add_node("handle_sku", handle_sku_selection_node)
    workflow.add_node("view_cart", view_cart_node)
    workflow.add_node("confirm_cart", confirm_cart_node)
    workflow.add_node("handle_cart", handle_cart_confirmation_node)
    workflow.add_node("collect_address", collect_address_node)
    workflow.add_node("handle_address", handle_address_node)
    workflow.add_node("confirm_order", confirm_order_node)
    workflow.add_node("handle_order", handle_order_confirmation_node)
    workflow.add_node("place_order", place_order_node)
    
    # 设置入口
    workflow.set_entry_point("search")
    
    # 添加边
    workflow.add_edge("search", "confirm_sku")
    workflow.add_edge("confirm_sku", "handle_sku")
    workflow.add_conditional_edges(
        "handle_sku",
        route_after_sku_handling,
        {
            "confirm_sku": "confirm_sku",
            "view_cart": "view_cart",
            "search": "search"
        }
    )
    
    workflow.add_edge("view_cart", "confirm_cart")
    workflow.add_edge("confirm_cart", "handle_cart")
    workflow.add_conditional_edges(
        "handle_cart",
        route_after_cart_handling,
        {
            "confirm_cart": "confirm_cart",
            "collect_address": "collect_address",
            "search": "search"
        }
    )
    
    workflow.add_edge("collect_address", "handle_address")
    workflow.add_conditional_edges(
        "handle_address",
        route_after_address_handling,
        {
            "collect_address": "collect_address",
            "confirm_order": "confirm_order"
        }
    )
    
    workflow.add_edge("confirm_order", "handle_order")
    workflow.add_conditional_edges(
        "handle_order",
        route_after_order_confirmation,
        {
            "confirm_order": "confirm_order",
            "place_order": "place_order",
            "collect_address": "collect_address"
        }
    )
    
    workflow.add_conditional_edges(
        "place_order",
        route_after_place_order,
        {"end": END}
    )
    
    # 使用内存checkpointer，在关键节点前interrupt
    memory = MemorySaver()
    return workflow.compile(
        checkpointer=memory,
        interrupt_before=["confirm_sku", "confirm_cart", "collect_address", "confirm_order"]
    )
