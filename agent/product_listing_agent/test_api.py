#!/usr/bin/env python3
"""Test script for API client and tools."""

import os
import sys
from dotenv import load_dotenv

# Add parent directory to path
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from product_listing_agent.api_client import PmallAPIClient
from product_listing_agent.tools import initialize_tools, get_categories_tool, get_brands_tool


def test_api_client():
    """Test the API client."""
    print("="*60)
    print("测试 API 客户端")
    print("="*60)
    
    # Load environment variables
    load_dotenv()
    
    # Initialize client
    print("\n1. 初始化客户端...")
    client = PmallAPIClient()
    print(f"   API URL: {client.base_url}")
    print(f"   用户名: {client.username}")
    
    # Test login
    print("\n2. 测试登录...")
    try:
        result = client.login()
        print(f"   ✅ 登录成功！")
        print(f"   用户ID: {result['user']['id']}")
        print(f"   用户名: {result['user']['username']}")
        print(f"   Token: {result['token'][:20]}...")
    except Exception as e:
        print(f"   ❌ 登录失败: {e}")
        return False
    
    # Test get categories
    print("\n3. 测试获取分类...")
    try:
        categories = client.get_categories()
        print(f"   ✅ 获取成功！共 {len(categories)} 个分类")
        if categories:
            print(f"   示例: {categories[0]['name']} (ID: {categories[0]['id']})")
    except Exception as e:
        print(f"   ❌ 获取失败: {e}")
    
    # Test get brands
    print("\n4. 测试获取品牌...")
    try:
        result = client.get_brands(page=1, page_size=10)
        brands = result.get('brands', [])
        print(f"   ✅ 获取成功！共 {result.get('total', 0)} 个品牌")
        if brands:
            print(f"   示例: {brands[0]['name']} (ID: {brands[0]['id']})")
    except Exception as e:
        print(f"   ❌ 获取失败: {e}")
    
    # Initialize tools
    print("\n5. 初始化工具...")
    try:
        initialize_tools(client)
        print("   ✅ 工具初始化成功！")
    except Exception as e:
        print(f"   ❌ 工具初始化失败: {e}")
        return False
    
    # Test tools
    print("\n6. 测试工具调用...")
    try:
        result = get_categories_tool.invoke({})
        print("   ✅ get_categories_tool 调用成功！")
        print(f"   返回数据长度: {len(result)} 字符")
    except Exception as e:
        print(f"   ❌ 工具调用失败: {e}")
    
    try:
        result = get_brands_tool.invoke({"page": 1, "page_size": 5})
        print("   ✅ get_brands_tool 调用成功！")
        print(f"   返回数据长度: {len(result)} 字符")
    except Exception as e:
        print(f"   ❌ 工具调用失败: {e}")
    
    print("\n" + "="*60)
    print("✅ 所有测试完成！")
    print("="*60)
    return True


def test_create_product():
    """Test creating a sample product."""
    print("\n" + "="*60)
    print("测试创建商品")
    print("="*60)
    
    load_dotenv()
    
    client = PmallAPIClient()
    client.login()
    
    # Sample product data
    sample_product = {
        "spu": {
            "brand_id": 1,
            "category_id": 1,
            "name": "测试商品-iPhone 15 Pro",
            "sub_title": "256GB 钛金属",
            "main_image": "https://example.com/iphone15pro.jpg",
            "sort": 0,
            "service_bits": 0
        },
        "skus": [
            {
                "sku_code": "TEST-IP15P-256-BLACK",
                "name": "iPhone 15 Pro 256GB 黑色",
                "sub_title": "黑色钛金属",
                "main_image": "https://example.com/iphone15pro-black.jpg",
                "price": "8999.00",
                "market_price": "9999.00",
                "stock": 100,
                "sku_spec_data": '{"color": "黑色", "storage": "256GB"}'
            }
        ],
        "detail": {
            "description": "这是一个测试商品",
            "images": ["https://example.com/detail1.jpg"],
            "videos": [],
            "market_tag_json": "{}",
            "tech_tag_json": "{}"
        }
    }
    
    print("\n准备创建测试商品...")
    print(f"商品名称: {sample_product['spu']['name']}")
    
    confirm = input("\n确认创建测试商品? (y/n): ")
    if confirm.lower() != 'y':
        print("已取消")
        return
    
    try:
        result = client.create_product(
            spu=sample_product['spu'],
            skus=sample_product['skus'],
            detail=sample_product['detail']
        )
        print(f"\n✅ 创建成功！")
        print(f"SPU ID: {result.get('spu_id')}")
        print(f"消息: {result.get('message')}")
    except Exception as e:
        print(f"\n❌ 创建失败: {e}")


if __name__ == "__main__":
    import argparse
    
    parser = argparse.ArgumentParser(description="测试 API 客户端和工具")
    parser.add_argument(
        "--create",
        action="store_true",
        help="测试创建商品（会实际创建数据）"
    )
    
    args = parser.parse_args()
    
    # Run basic tests
    success = test_api_client()
    
    # Optionally test product creation
    if args.create and success:
        test_create_product()
