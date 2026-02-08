"""Main entry point for Auto Order Agent

提供命令行交互界面
"""

import os
import sys
from dotenv import load_dotenv
from langchain_core.messages import HumanMessage

from .api_client import PmallOrderAPIClient
from .tools import initialize_tools
from .state import create_initial_state
from .graph import build_graph

# 加载环境变量
load_dotenv()


def print_banner():
    """打印欢迎横幅"""
    banner = """
╔══════════════════════════════════════════════════════════╗
║                                                          ║
║        🛒  自动下单 Agent  🛒                             ║
║                                                          ║
║  智能购物助手：搜索商品 → 加入购物车 → 一键下单            ║
║                                                          ║
╚══════════════════════════════════════════════════════════╝
    """
    print(banner)


def print_help():
    """打印帮助信息"""
    help_text = """
📖 使用说明：

1️⃣ 搜索商品：
   "我想买红米手机"
   "找一个2000元左右的手机"

2️⃣ 查看购物车：
   "查看购物车"
   "购物车"

3️⃣ 结算下单：
   "去结算"
   "下单"

4️⃣ 其他命令：
   exit/quit/退出 - 退出程序
   help/帮助 - 显示帮助

💡 提示：Agent会引导你完成整个购物流程！
    """
    print(help_text)


def run_interactive():
    """运行交互式界面"""
    print_banner()
    
    # 检查环境变量
    required_env_vars = ["OPENAI_API_KEY", "OPENAI_BASE_URL"]
    missing_vars = [var for var in required_env_vars if not os.getenv(var)]
    
    if missing_vars:
        print(f"❌ 缺少环境变量: {', '.join(missing_vars)}")
        print("\n请设置以下环境变量：")
        print("  export OPENAI_API_KEY='your-api-key'")
        print("  export OPENAI_BASE_URL='https://api.openai.com/v1'")
        print("  export OPENAI_MODEL='gpt-4'  # 可选")
        return
    
    # 初始化API客户端
    print("🔐 正在登录...")
    api_client = PmallOrderAPIClient(base_url="http://localhost:8080")
    
    try:
        result = api_client.login("piao", "123456")
        if not api_client.token:
            print(f"❌ 登录失败: {result}")
            return
        print("✅ 登录成功！\n")
    except Exception as e:
        print(f"❌ 连接API失败: {e}")
        print("请确保API服务运行在 http://localhost:8080")
        return
    
    # 初始化工具
    initialize_tools(api_client)
    
    # 构建graph
    graph = build_graph()
    
    print("✨ 准备就绪！请告诉我您想买什么：\n")
    
    # 线程配置
    thread_id = 1
    
    # 主循环
    while True:
        try:
            user_input = input("👤 你: ").strip()
            
            if not user_input:
                continue
            
            if user_input.lower() in ['exit', 'quit', '退出']:
                print("\n👋 感谢使用，再见！")
                break
            
            if user_input.lower() in ['help', '帮助']:
                print_help()
                continue
            
            # 创建新会话
            config = {"configurable": {"thread_id": str(thread_id)}}
            initial_state = create_initial_state(user_input)
            
            print("\n🤖 Agent: 正在处理...\n")
            
            # 会话循环
            session_active = True
            while session_active:
                # 执行graph直到interrupt或结束
                for event in graph.stream(initial_state, config, stream_mode="values"):
                    if "messages" in event:
                        last_message = event["messages"][-1]
                        if hasattr(last_message, 'content') and last_message.content:
                            if not isinstance(last_message, HumanMessage):
                                print(f"🤖 Agent: {last_message.content}\n")
                    
                    # 检查是否结束
                    if event.get("next_step") == "end":
                        thread_id += 1
                        session_active = False
                        print("\n" + "="*60)
                        print("会话结束。继续购物或输入 'exit' 退出。")
                        print("="*60 + "\n")
                        break
                
                # 检查是否在interrupt点
                if session_active:
                    snapshot = graph.get_state(config)
                    
                    # 检测当前interrupt的节点
                    interrupt_nodes = ["confirm_sku", "confirm_cart", "collect_address", "confirm_order"]
                    current_interrupt = None
                    
                    if snapshot.next:
                        for node in interrupt_nodes:
                            if node in snapshot.next:
                                current_interrupt = node
                                break
                    
                    if current_interrupt:
                        # 获取用户输入
                        user_response = input("👤 你: ").strip()
                        
                        if user_response.lower() in ['exit', 'quit', '退出']:
                            print("\n👋 感谢使用，再见！")
                            return
                        
                        # 更新状态
                        graph.update_state(
                            config,
                            {"messages": [HumanMessage(content=user_response)]},
                            as_node=current_interrupt
                        )
                        
                        # 继续执行
                        initial_state = None
                    else:
                        # 没有interrupt点，会话结束
                        session_active = False
        
        except KeyboardInterrupt:
            print("\n\n👋 感谢使用，再见！")
            break
        except Exception as e:
            print(f"\n❌ 发生错误: {e}")
            import traceback
            traceback.print_exc()
            print("\n会话已重置，请重新开始。\n")


def main():
    """主入口"""
    run_interactive()


if __name__ == "__main__":
    main()
