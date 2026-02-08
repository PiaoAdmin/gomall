"""Main entry point for the product listing agent.

This module provides a command-line interface for the product listing agent.
"""

import os
import sys
from typing import Optional
from langchain_core.messages import HumanMessage

from .api_client import PmallAPIClient
from .tools import initialize_tools
from .state import create_initial_state
from .graph import build_graph


def print_banner():
    """Print welcome banner."""
    banner = """
╔══════════════════════════════════════════════════════════╗
║                                                          ║
║        🛍️  商品自助上架 Agent  🛍️                         ║
║                                                          ║
║  支持自然语言输入，智能补全商品信息，一键上架！            ║
║                                                          ║
╚══════════════════════════════════════════════════════════╝
    """
    print(banner)


def print_help():
    """Print help information."""
    help_text = """
使用说明：
---------
1. 描述你想要上架的商品，例如：
   "帮我添加一个iPhone 15 Pro Max，256GB，价格8999元，库存100"
   
2. Agent 会自动补全商品信息并展示给你确认

3. 你可以：
   - 输入「是」、「确认」等同意创建
   - 输入修改意见，如"价格改成8888"
   
4. 确认后 Agent 会自动创建商品

5. 输入 'exit' 或 'quit' 退出程序

环境变量配置：
-------------
- PMALL_API_URL: API服务地址（默认: http://localhost:8888）
- PMALL_USERNAME: 登录用户名（默认: piao）
- PMALL_PASSWORD: 登录密码（默认: 123456）
- OPENAI_API_KEY: OpenAI API密钥（必需）
- OPENAI_API_BASE: OpenAI API地址（可选）
- OPENAI_MODEL: 模型名称（默认: qwen-plus）
    """
    print(help_text)


def run_interactive():
    """Run the agent in interactive mode."""
    print_banner()
    print_help()
    
    # Initialize API client
    print("\n正在初始化...")
    try:
        api_client = PmallAPIClient()
        print(f"🔗 连接到: {api_client.base_url}")
        
        # Login
        print(f"🔐 登录用户: {api_client.username}")
        login_result = api_client.login()
        print(f"✅ 登录成功！用户ID: {login_result['user']['id']}\n")
        
        # Initialize tools with the API client
        initialize_tools(api_client)
        
        # Build the graph
        graph = build_graph()
        
    except Exception as e:
        print(f"❌ 初始化失败: {e}")
        print("\n请检查：")
        print("1. API服务是否启动（默认 http://localhost:8080）")
        print("2. 环境变量是否正确配置")
        print("3. 用户名密码是否正确")
        sys.exit(1)
    
    print("✨ 准备就绪！请描述你要上架的商品：\n")
    
    # Thread configuration for maintaining conversation state
    thread_id = 1
    
    # Main interaction loop
    while True:
        try:
            user_input = input("👤 你: ").strip()
            
            if not user_input:
                continue
            
            if user_input.lower() in ['exit', 'quit', '退出']:
                print("\n👋 再见！")
                break
            
            if user_input.lower() in ['help', '帮助']:
                print_help()
                continue
            
            # Create new thread for this product listing session
            config = {"configurable": {"thread_id": str(thread_id)}}
            
            # Create initial state with user input
            print("\n🤖 Agent: 正在分析商品信息...\n")
            initial_state = create_initial_state(user_input)
            
            # Stream the graph execution
            session_active = True
            while session_active:
                # Run graph until interrupt or completion
                for event in graph.stream(initial_state, config, stream_mode="values"):
                    # Print AI messages
                    if "messages" in event:
                        last_message = event["messages"][-1]
                        if hasattr(last_message, 'content') and last_message.content:
                            # Skip user messages (they're already printed)
                            if not isinstance(last_message, HumanMessage):
                                print(f"🤖 Agent: {last_message.content}\n")
                    
                    # Check if workflow ended
                    if event.get("next_step") == "end":
                        thread_id += 1  # Increment for next session
                        session_active = False
                        print("\n" + "="*60)
                        print("会话结束。请描述下一个要上架的商品，或输入 'exit' 退出。")
                        print("="*60 + "\n")
                        break
                
                # If session not ended, check if we're at an interrupt point
                if session_active:
                    snapshot = graph.get_state(config)
                    
                    # Check if we're interrupted (at confirm node)
                    if snapshot.next and "confirm" in snapshot.next:
                        # Get user response
                        user_response = input("👤 你: ").strip()
                        
                        if user_response.lower() in ['exit', 'quit', '退出']:
                            print("\n👋 再见！")
                            return
                        
                        # Update state with user response
                        graph.update_state(
                            config,
                            {"messages": [HumanMessage(content=user_response)]},
                            as_node="confirm"
                        )
                        
                        # Continue from checkpoint (set initial_state to None)
                        initial_state = None
                    else:
                        # No interrupt, session must have ended
                        session_active = False
        
        except KeyboardInterrupt:
            print("\n\n👋 再见！")
            break
        except Exception as e:
            print(f"\n❌ 发生错误: {e}")
            import traceback
            traceback.print_exc()
            print("\n会话已重置，请重新开始。\n")


def main():
    """Main entry point."""
    run_interactive()


if __name__ == "__main__":
    main()
