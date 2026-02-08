"""LangChain tools for product listing operations.

This module provides LangChain-compatible tools that wrap the pmall API client
for use by the agent.
"""

from typing import Optional, Dict, Any, List
from langchain.tools import tool
from .api_client import PmallAPIClient

# Global API client instance
_api_client: Optional[PmallAPIClient] = None


def initialize_tools(api_client: PmallAPIClient):
    """Initialize the tools with an API client instance.
    
    Args:
        api_client: Initialized PmallAPIClient instance
    """
    global _api_client
    _api_client = api_client


def get_api_client() -> PmallAPIClient:
    """Get the global API client instance.
    
    Returns:
        The initialized API client
        
    Raises:
        RuntimeError: If tools haven't been initialized
    """
    if _api_client is None:
        raise RuntimeError("Tools not initialized. Call initialize_tools() first.")
    return _api_client


@tool
def get_categories_tool(parent_id: Optional[int] = None) -> str:
    """获取商品分类列表。
    
    Args:
        parent_id: 父分类ID，不传或传None则获取根分类
        
    Returns:
        分类列表的JSON字符串，包含id、name、level等信息
    """
    import json
    client = get_api_client()
    categories = client.get_categories(parent_id)
    return json.dumps(categories, ensure_ascii=False, indent=2)


@tool
def get_brands_tool(page: int = 1, page_size: int = 50) -> str:
    """获取品牌列表。
    
    Args:
        page: 页码，默认1
        page_size: 每页数量，默认50
        
    Returns:
        品牌列表的JSON字符串，包含id、name、logo等信息
    """
    import json
    client = get_api_client()
    result = client.get_brands(page, page_size)
    return json.dumps(result, ensure_ascii=False, indent=2)


@tool
def create_product_tool(product_data: str) -> str:
    """创建商品。必须提供完整的商品信息。
    
    Args:
        product_data: JSON字符串格式的商品数据，必须包含以下结构：
            {
                "spu": {
                    "brand_id": int,           # 品牌ID（必需）
                    "category_id": int,        # 分类ID（必需）
                    "name": str,               # 商品名称（必需）
                    "sub_title": str,          # 副标题
                    "main_image": str,         # 主图URL
                    "sort": int,               # 排序，默认0
                    "service_bits": int        # 服务标识，默认0
                },
                "skus": [                      # SKU列表（至少1个）
                    {
                        "sku_code": str,       # SKU编码（必需）
                        "name": str,           # SKU名称（必需）
                        "sub_title": str,      # 副标题
                        "main_image": str,     # 主图URL
                        "price": str,          # 价格，如"99.99"（必需）
                        "market_price": str,   # 市场价，如"199.99"
                        "stock": int,          # 库存（必需）
                        "sku_spec_data": str   # 规格JSON字符串
                    }
                ],
                "detail": {                    # 详情（可选）
                    "description": str,        # 富文本描述
                    "images": [str],           # 详情图片URL列表
                    "videos": [str],           # 视频URL列表
                    "market_tag_json": str,    # 营销标签JSON
                    "tech_tag_json": str       # 技术参数JSON
                }
            }
    
    Returns:
        创建结果的JSON字符串，包含spu_id和message
    """
    import json
    client = get_api_client()
    
    try:
        data = json.loads(product_data)
        
        # 验证必需字段
        if "spu" not in data:
            return json.dumps({"error": "缺少spu字段"}, ensure_ascii=False)
        if "skus" not in data or not data["skus"]:
            return json.dumps({"error": "skus字段不能为空，至少需要1个SKU"}, ensure_ascii=False)
        
        spu = data["spu"]
        skus = data["skus"]
        detail = data.get("detail")
        
        # 验证SPU必需字段
        required_spu_fields = ["brand_id", "category_id", "name"]
        missing_spu = [f for f in required_spu_fields if f not in spu or not spu[f]]
        if missing_spu:
            return json.dumps({"error": f"SPU缺少必需字段: {', '.join(missing_spu)}"}, ensure_ascii=False)
        
        # 验证SKU必需字段
        for i, sku in enumerate(skus):
            required_sku_fields = ["sku_code", "name", "price", "stock"]
            missing_sku = [f for f in required_sku_fields if f not in sku or (sku[f] is None and f != "stock")]
            if missing_sku:
                return json.dumps({"error": f"SKU[{i}]缺少必需字段: {', '.join(missing_sku)}"}, ensure_ascii=False)
        
        result = client.create_product(spu, skus, detail)
        return json.dumps(result, ensure_ascii=False, indent=2)
    
    except json.JSONDecodeError as e:
        return json.dumps({"error": f"JSON解析错误: {str(e)}"}, ensure_ascii=False)
    except Exception as e:
        return json.dumps({"error": f"创建商品失败: {str(e)}"}, ensure_ascii=False)


# Export all tools
ALL_TOOLS = [
    get_categories_tool,
    get_brands_tool,
    create_product_tool,
]
