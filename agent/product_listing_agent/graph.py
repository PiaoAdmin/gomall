"""LangGraph workflow for product listing agent.

This module implements the state machine for the product listing workflow.
"""

import json
import os
from typing import Literal
from langgraph.graph import StateGraph, END
from langchain_openai import ChatOpenAI
from langchain_core.messages import HumanMessage, AIMessage, SystemMessage

from .state import (
    ProductListingState,
    COMPLETE_INFO_SYSTEM_PROMPT,
    VALIDATION_SYSTEM_PROMPT,
    ERROR_RETRY_SYSTEM_PROMPT,
    format_product_data_for_display,
)
from .tools import ALL_TOOLS


# Initialize LLM
def get_llm():
    """Get the LLM instance."""
    # Support both OpenAI-compatible APIs and Dashscope
    api_base = os.getenv("OPENAI_API_BASE") or os.getenv("OPENAI_BASE_URL")
    api_key = os.getenv("OPENAI_API_KEY")
    model = os.getenv("OPENAI_MODEL", "qwen-plus")
    
    if not api_key:
        raise ValueError("OPENAI_API_KEY environment variable is required")
    
    return ChatOpenAI(
        model=model,
        openai_api_key=api_key,
        openai_api_base=api_base,
        temperature=0.7,
    )


def complete_product_info_node(state: ProductListingState) -> ProductListingState:
    """Node: LLM completes the product information based on user input.
    
    This node uses LLM with tools to:
    1. Extract product info from user input
    2. Query categories and brands if needed
    3. Complete the full product data structure
    """
    llm = get_llm()
    llm_with_tools = llm.bind_tools(ALL_TOOLS)
    
    messages = [
        SystemMessage(content=COMPLETE_INFO_SYSTEM_PROMPT),
        *state["messages"]
    ]
    
    # Let LLM process with tools
    response = llm_with_tools.invoke(messages)
    
    # If LLM wants to use tools, execute them
    if response.tool_calls:
        # Execute tool calls
        from langchain_core.messages import ToolMessage
        
        messages.append(response)
        
        for tool_call in response.tool_calls:
            # Find and execute the tool
            tool_name = tool_call["name"]
            tool_args = tool_call["args"]
            
            for tool in ALL_TOOLS:
                if tool.name == tool_name:
                    try:
                        result = tool.invoke(tool_args)
                        messages.append(ToolMessage(
                            content=result,
                            tool_call_id=tool_call["id"]
                        ))
                    except Exception as e:
                        messages.append(ToolMessage(
                            content=f"Error: {str(e)}",
                            tool_call_id=tool_call["id"]
                        ))
                    break
        
        # Get final response after tool execution
        response = llm_with_tools.invoke(messages)
    
    # Try to extract JSON product data from response
    content = response.content
    product_data = None
    
    try:
        # Try to find JSON in the response (look for ```json blocks first)
        if "```json" in content:
            start_idx = content.find("```json") + 7
            end_idx = content.find("```", start_idx)
            json_str = content[start_idx:end_idx].strip()
        elif "```" in content:
            start_idx = content.find("```") + 3
            end_idx = content.find("```", start_idx)
            json_str = content[start_idx:end_idx].strip()
        else:
            # Find raw JSON
            start_idx = content.find('{')
            end_idx = content.rfind('}') + 1
            if start_idx >= 0 and end_idx > start_idx:
                json_str = content[start_idx:end_idx]
            else:
                json_str = None
        
        if json_str:
            product_data = json.loads(json_str)
            
            # Validate required structure
            if "spu" not in product_data or "skus" not in product_data:
                product_data = None
    except Exception as e:
        # If parsing fails, try to ask LLM to fix it
        print(f"JSON parsing error: {e}")
        pass
    
    # Update state
    # Don't add the raw LLM response to messages, only add the formatted confirmation message
    # This avoids duplicate message issues
    new_messages = state["messages"]
    
    if product_data:
        # Show formatted product data to user
        formatted = format_product_data_for_display(product_data)
        confirmation_msg = f"{formatted}\n\n请确认以上商品信息是否正确？\n- 输入「是」或「确认」继续创建\n- 输入修改意见进行调整"
        new_messages = new_messages + [AIMessage(content=confirmation_msg)]
        next_step = "confirm"
    else:
        # Failed to extract proper JSON - force LLM to generate data by emphasizing requirement
        # Add a strong directive to generate data
        retry_messages = [
            SystemMessage(content=COMPLETE_INFO_SYSTEM_PROMPT + "\n\n**CRITICAL: You MUST generate complete JSON data structure even with minimal information! Use reasonable defaults for missing fields. This is REQUIRED, not optional!**"),
            *state["messages"],
            AIMessage(content="我明白了，即使信息有限，我也必须生成完整的JSON数据。让我基于已有信息生成合理的商品数据。"),
        ]
        
        # Try one more time with stronger prompting
        retry_response = llm_with_tools.invoke(retry_messages)
        retry_content = retry_response.content
        
        # Try to extract JSON again
        try:
            if "```json" in retry_content:
                start_idx = retry_content.find("```json") + 7
                end_idx = retry_content.find("```", start_idx)
                json_str = retry_content[start_idx:end_idx].strip()
            elif "```" in retry_content:
                start_idx = retry_content.find("```") + 3
                end_idx = retry_content.find("```", start_idx)
                json_str = retry_content[start_idx:end_idx].strip()
            else:
                start_idx = retry_content.find('{')
                end_idx = retry_content.rfind('}') + 1
                if start_idx >= 0 and end_idx > start_idx:
                    json_str = retry_content[start_idx:end_idx]
                else:
                    json_str = None
            
            if json_str:
                product_data = json.loads(json_str)
                if "spu" in product_data and "skus" in product_data:
                    # Success! Show the data
                    formatted = format_product_data_for_display(product_data)
                    confirmation_msg = f"{formatted}\n\n请确认以上商品信息是否正确？\n- 输入「是」或「确认」继续创建\n- 输入修改意见进行调整"
                    new_messages = new_messages + [AIMessage(content=confirmation_msg)]
                    next_step = "confirm"
                else:
                    # Still no valid data, ask user for more info
                    new_messages = new_messages + [AIMessage(content="抱歉，信息太少，无法生成商品数据。\n\n请提供更多信息，例如：\n- 具体型号和配置（如'红米K70 12GB+256GB'）\n- 期望价格\n- 其他详细信息")]
                    product_data = None
                    next_step = "confirm"
        except:
            # Final fallback
            new_messages = new_messages + [AIMessage(content="抱歉，信息太少，无法生成商品数据。\n\n请提供更多信息，例如：\n- 具体型号和配置（如'红米K70 12GB+256GB'）\n- 期望价格\n- 其他详细信息")]
            next_step = "confirm"
    
    return {
        **state,
        "messages": new_messages,
        "product_data": product_data,
        "next_step": next_step
    }


def user_confirmation_node(state: ProductListingState) -> ProductListingState:
    """Node: Wait for user confirmation or modification request.
    
    This is a human-in-the-loop node. The workflow will pause here
    via interrupt_before mechanism and wait for user input.
    """
    # This node does nothing - it's just a placeholder for the interrupt
    # The graph will pause BEFORE this node executes
    return state


def validate_user_input_node(state: ProductListingState) -> ProductListingState:
    """Node: Validate user's confirmation or modification request.
    
    This node analyzes the user's response to determine:
    - If approved: proceed to creation
    - If rejected: update product data and ask for confirmation again
    """
    llm = get_llm()
    
    # Get the last user message
    last_message = None
    for msg in reversed(state["messages"]):
        if isinstance(msg, HumanMessage):
            last_message = msg.content
            break
    
    if not last_message:
        return {**state, "next_step": "confirm"}
    
    # Check for simple approval keywords (more flexible matching)
    approval_keywords = ["是", "确认", "ok", "yes", "好", "可以", "没问题", "同意", "行"]
    user_input_normalized = last_message.lower().strip()
    
    # Remove common typos or extra characters
    if any(keyword in user_input_normalized for keyword in approval_keywords):
        return {
            **state,
            "validation_status": "approved",
            "messages": state["messages"] + [AIMessage(content="好的，开始创建商品...")],
            "next_step": "create"
        }
    
    # User wants to modify, use LLM to update product data
    messages = [
        SystemMessage(content=VALIDATION_SYSTEM_PROMPT),
        SystemMessage(content=f"当前商品数据：\n{json.dumps(state['product_data'], ensure_ascii=False, indent=2)}"),
        HumanMessage(content=f"用户反馈：{last_message}")
    ]
    
    response = llm.invoke(messages)
    
    try:
        # Extract updated product data
        content = response.content
        start_idx = content.find('{')
        end_idx = content.rfind('}') + 1
        if start_idx >= 0 and end_idx > start_idx:
            json_str = content[start_idx:end_idx]
            result = json.loads(json_str)
            
            if result.get("action") == "approved":
                return {
                    **state,
                    "validation_status": "approved",
                    "messages": state["messages"] + [AIMessage(content="好的，开始创建商品...")],
                    "next_step": "create"
                }
            else:
                # Update product data and ask for confirmation again
                updated_data = result.get("data", state["product_data"])
                formatted = format_product_data_for_display(updated_data)
                confirmation_msg = f"已根据您的要求更新：\n\n{formatted}\n\n请确认以上商品信息是否正确？\n- 输入「是」或「确认」继续创建\n- 输入修改意见进行调整"
                
                return {
                    **state,
                    "product_data": updated_data,
                    "messages": state["messages"] + [AIMessage(content=confirmation_msg)],
                    "validation_status": "pending",
                    "next_step": "confirm"  # Go back to confirm for another round
                }
    except Exception as e:
        # If parsing fails, ask user to clarify
        return {
            **state,
            "messages": state["messages"] + [AIMessage(content=f"抱歉，我没有理解您的修改意见。请明确指出需要修改的字段和新的值。\n错误：{str(e)}")],
            "next_step": "confirm"  # Go back to confirm
        }
    
    return {**state, "messages": state["messages"], "next_step": "confirm"}


def create_product_node(state: ProductListingState) -> ProductListingState:
    """Node: Call API to create the product.
    
    This node uses the create_product_tool to actually create the product.
    """
    from .tools import create_product_tool
    
    if not state["product_data"]:
        return {
            **state,
            "error_message": "商品数据为空，无法创建",
            "next_step": "end"
        }
    
    try:
        # Convert product data to JSON string for the tool
        product_json = json.dumps(state["product_data"], ensure_ascii=False)
        
        # Call the create product tool
        result_str = create_product_tool.invoke({"product_data": product_json})
        result = json.loads(result_str)
        
        if "error" in result:
            # Creation failed
            return {
                **state,
                "error_message": result["error"],
                "retry_count": state["retry_count"] + 1,
                "messages": [AIMessage(content=f"❌ 创建失败：{result['error']}")],
                "next_step": "retry"
            }
        else:
            # Success
            success_msg = f"✅ 商品创建成功！\nSPU ID: {result.get('spu_id')}\n{result.get('message', '')}"
            return {
                **state,
                "messages": [AIMessage(content=success_msg)],
                "next_step": "end"
            }
    
    except Exception as e:
        return {
            **state,
            "error_message": str(e),
            "retry_count": state["retry_count"] + 1,
            "messages": [AIMessage(content=f"❌ 创建失败：{str(e)}")],
            "next_step": "retry"
        }


def retry_with_fix_node(state: ProductListingState) -> ProductListingState:
    """Node: Retry creation after fixing errors.
    
    This node uses LLM to analyze the error and fix the product data.
    """
    if state["retry_count"] >= 3:
        return {
            **state,
            "messages": [AIMessage(content=f"❌ 已重试3次仍然失败，请检查数据或联系管理员。\n最后错误：{state['error_message']}")],
            "next_step": "end"
        }
    
    llm = get_llm()
    
    prompt = ERROR_RETRY_SYSTEM_PROMPT.format(error_message=state["error_message"])
    messages = [
        SystemMessage(content=prompt),
        SystemMessage(content=f"当前商品数据：\n{json.dumps(state['product_data'], ensure_ascii=False, indent=2)}")
    ]
    
    response = llm.invoke(messages)
    
    try:
        # Extract fixed product data
        content = response.content
        start_idx = content.find('{')
        end_idx = content.rfind('}') + 1
        if start_idx >= 0 and end_idx > start_idx:
            json_str = content[start_idx:end_idx]
            fixed_data = json.loads(json_str)
            
            return {
                **state,
                "product_data": fixed_data,
                "messages": [AIMessage(content=f"🔧 已修复数据问题，重试第{state['retry_count']}次...")],
                "next_step": "create"
            }
    except Exception as e:
        return {
            **state,
            "messages": [AIMessage(content=f"❌ 无法自动修复错误：{str(e)}")],
            "next_step": "end"
        }
    
    return {**state, "next_step": "end"}


def route_after_validation(state: ProductListingState) -> str:
    """Router: Determine next step after validation."""
    if state["next_step"] == "create":
        return "create"
    # Loop back to confirm for more modifications
    return "confirm"


def route_after_creation(state: ProductListingState) -> str:
    """Router: Determine next step after creation attempt."""
    if state["next_step"] == "retry":
        return "retry"
    return "end"


def route_after_retry(state: ProductListingState) -> str:
    """Router: Determine next step after retry."""
    if state["next_step"] == "create":
        return "create"
    return "end"


def build_graph():
    """Build the LangGraph workflow.
    
    Returns:
        Compiled state graph with memory checkpointer
    """
    from langgraph.checkpoint.memory import MemorySaver
    
    workflow = StateGraph(ProductListingState)
    
    # Add nodes
    workflow.add_node("complete_info", complete_product_info_node)
    workflow.add_node("confirm", user_confirmation_node)
    workflow.add_node("validate", validate_user_input_node)
    workflow.add_node("create", create_product_node)
    workflow.add_node("retry", retry_with_fix_node)
    
    # Set entry point
    workflow.set_entry_point("complete_info")
    
    # Add edges
    workflow.add_edge("complete_info", "confirm")
    workflow.add_edge("confirm", "validate")  # After user input, go to validate
    
    workflow.add_conditional_edges(
        "validate",
        route_after_validation,
        {
            "create": "create",
            "confirm": "confirm"  # Loop back for more modifications
        }
    )
    workflow.add_conditional_edges(
        "create",
        route_after_creation,
        {
            "retry": "retry",
            "end": END
        }
    )
    workflow.add_conditional_edges(
        "retry",
        route_after_retry,
        {
            "create": "create",
            "end": END
        }
    )
    
    # Use memory checkpointer and interrupt before confirm node
    memory = MemorySaver()
    return workflow.compile(checkpointer=memory, interrupt_before=["confirm"])
