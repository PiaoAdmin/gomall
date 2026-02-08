#!/usr/bin/env python3
"""简单的非交互式测试脚本，验证 Agent 的基本功能"""

import os
import sys
from dotenv import load_dotenv

# Add parent directory to path
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from product_listing_agent.api_client import PmallAPIClient
from product_listing_agent.tools import initialize_tools, get_categories_tool, get_brands_tool
from product_listing_agent.state import create_initial_state, format_product_data_for_display
from product_listing_agent.graph import get_llm

def test_llm_connection():
    """测试 LLM 连接"""
    print("\n" + "="*60)
    print("测试 LLM 连接")
    print("="*60)
    
    load_dotenv()
    
    try:
        llm = get_llm()
        response = llm.invoke("你好，请回复'测试成功'")
        print(f"✅ LLM 连接成功！")
        print(f"响应: {response.content[:100]}...")
        return True
    except Exception as e:
        print(f"❌ LLM 连接失败: {e}")
        return False


def test_state_creation():
    """测试状态创建"""
    print("\n" + "="*60)
    print("测试状态创建")
    print("="*60)
    
    try:
        user_input = "添加一个测试商品"
        state = create_initial_state(user_input)
        
        print(f"✅ 状态创建成功！")
        print(f"   user_input: {state['user_input']}")
        print(f"   validation_status: {state['validation_status']}")
        print(f"   next_step: {state['next_step']}")
        print(f"   messages count: {len(state['messages'])}")
        return True
    except Exception as e:
        print(f"❌ 状态创建失败: {e}")
        return False


def test_product_data_formatting():
    """测试商品数据格式化"""
    print("\n" + "="*60)
    print("测试商品数据格式化")
    print("="*60)
    
    try:
        sample_data = {
            "spu": {
                "brand_id": 1,
                "category_id": 1,
                "name": "测试商品",
                "sub_title": "这是一个测试"
            },
            "skus": [
                {
                    "sku_code": "TEST-001",
                    "name": "测试SKU",
                    "price": "99.99",
                    "stock": 100
                }
            ]
        }
        
        formatted = format_product_data_for_display(sample_data)
        print("✅ 数据格式化成功！")
        print("\n" + formatted[:300] + "...")
        return True
    except Exception as e:
        print(f"❌ 数据格式化失败: {e}")
        return False


def main():
    """运行所有测试"""
    print("\n🧪 开始运行 Agent 功能测试...\n")
    
    results = []
    
    # 测试 LLM 连接
    results.append(("LLM 连接", test_llm_connection()))
    
    # 测试状态创建
    results.append(("状态创建", test_state_creation()))
    
    # 测试数据格式化
    results.append(("数据格式化", test_product_data_formatting()))
    
    # 总结
    print("\n" + "="*60)
    print("测试总结")
    print("="*60)
    
    for name, success in results:
        status = "✅ 通过" if success else "❌ 失败"
        print(f"{name}: {status}")
    
    total = len(results)
    passed = sum(1 for _, success in results if success)
    
    print(f"\n总计: {passed}/{total} 通过")
    
    if passed == total:
        print("\n🎉 所有测试通过！Agent 已准备就绪！")
        return 0
    else:
        print(f"\n⚠️ {total - passed} 个测试失败")
        return 1


if __name__ == "__main__":
    sys.exit(main())
