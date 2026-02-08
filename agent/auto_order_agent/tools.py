"""LangChain Tools for Auto Order Agent

将API封装为LangChain工具供LLM调用
"""

from langchain_core.tools import tool
from typing import Optional, List
from .api_client import PmallOrderAPIClient


# 全局API客户端实例
_api_client: Optional[PmallOrderAPIClient] = None


def initialize_tools(api_client: PmallOrderAPIClient):
    """初始化工具，设置API客户端"""
    global _api_client
    _api_client = api_client


def _extract_data(result):
    """提取API响应中的数据
    
    API可能返回两种格式：
    1. {code: 200, data: {...}, message: "..."}
    2. 直接数据 {...}
    """
    if isinstance(result, dict):
        if "data" in result:
            return result["data"]
        return result
    return result


@tool
def search_products_tool(
    keyword: Optional[str] = None,
    category_id: Optional[int] = None,
    brand_id: Optional[int] = None,
    min_price: Optional[float] = None,
    max_price: Optional[float] = None,
    sort_by: str = "default",
    page_size: int = 5
) -> str:
    """⚠️ 必须使用此工具从pmall商城API搜索商品！禁止凭空推荐商品！
    
    用户提到任何商品关键词时（如"手机"、"小米"、"2000元左右"），必须立即调用此工具搜索。
    
    Args:
        keyword: 搜索关键词（如"手机"、"红米"、"iPhone"）- 用户输入的任何商品相关词
        category_id: 分类ID（可选）
        brand_id: 品牌ID（可选）
        min_price: 最低价格（可选，如用户说"2000元左右"可设置1900）
        max_price: 最高价格（可选，如用户说"2000元左右"可设置2100）
        sort_by: 排序方式，可选值：default/price_asc/price_desc/sale
        page_size: 返回数量，默认5个
    
    Returns:
        商品SKU列表的JSON字符串，包含sku_id、spu_id、name、price、stock等信息
    
    示例调用：
    - 用户："帮我找手机" → search_products_tool(keyword="手机")
    - 用户："小米" → search_products_tool(keyword="小米")
    - 用户："2000元左右的手机" → search_products_tool(keyword="手机", min_price=1900, max_price=2100)
    """
    if _api_client is None:
        return "Error: API client not initialized"
    
    try:
        result = _api_client.search_products(
            keyword=keyword,
            category_id=category_id,
            brand_id=brand_id,
            min_price=min_price,
            max_price=max_price,
            sort_by=sort_by,
            page=1,
            page_size=page_size
        )
        
        data = _extract_data(result)
        # API返回的是 {list: [...], total: N}，不是 {skus: [...]}
        skus = data.get("list", []) if isinstance(data, dict) else []
        
        if not skus:
            return "未找到匹配的商品"
        
        # 格式化输出
        import json
        return json.dumps(skus, ensure_ascii=False, indent=2)
    except Exception as e:
        return f"搜索出错: {str(e)}"


@tool
def get_product_detail_tool(spu_id: int) -> str:
    """获取商品详细信息，包括所有SKU规格
    
    Args:
        spu_id: 商品SPU ID
    
    Returns:
        商品详情的JSON字符串，包含所有SKU
    """
    if _api_client is None:
        return "Error: API client not initialized"
    
    try:
        result = _api_client.get_product_detail(spu_id)
        data = _extract_data(result)
        
        import json
        return json.dumps(data, ensure_ascii=False, indent=2)
    except Exception as e:
        return f"获取详情出错: {str(e)}"


@tool
def view_cart_tool() -> str:
    """查看购物车内容
    
    Returns:
        购物车详情的JSON字符串
    """
    if _api_client is None:
        return "Error: API client not initialized"
    
    try:
        result = _api_client.get_cart()
        data = _extract_data(result)
        
        # API返回: {items: [...], total_price: ...}
        items = data.get("items", []) if isinstance(data, dict) else []
        if not items:
            return "购物车为空"
        
        import json
        return json.dumps(data, ensure_ascii=False, indent=2)
    except Exception as e:
        return f"获取购物车出错: {str(e)}"


@tool
def add_to_cart_tool(sku_id: int, quantity: int = 1) -> str:
    """添加商品到购物车
    
    Args:
        sku_id: 商品SKU ID
        quantity: 数量，默认1
    
    Returns:
        操作结果信息
    """
    if _api_client is None:
        return "Error: API client not initialized"
    
    try:
        result = _api_client.add_to_cart(sku_id, quantity)
        data = _extract_data(result)
        
        # API返回格式可能是 {item: {...}} 或直接 {...}
        if isinstance(data, dict):
            # 可能是item字段，也可能直接就是商品信息
            cart_item = data.get("item", data)
            sku_name = cart_item.get("sku_name", cart_item.get("name", f"SKU {sku_id}"))
            return f"✅ 已添加到购物车: {sku_name} x {quantity}"
        return "✅ 商品已添加到购物车"
    except Exception as e:
        return f"❌ 添加出错: {str(e)}"


@tool
def remove_from_cart_tool(sku_ids: List[int]) -> str:
    """从购物车移除商品
    
    Args:
        sku_ids: 要移除的SKU ID列表
    
    Returns:
        操作结果信息
    """
    if _api_client is None:
        return "Error: API client not initialized"
    
    try:
        result = _api_client.remove_from_cart(sku_ids)
        return f"已从购物车移除 {len(sku_ids)} 个商品"
    except Exception as e:
        return f"移除出错: {str(e)}"


@tool
def place_order_tool(email: str, name: str, street_address: str, city: str, zip_code: int) -> str:
    """下单（从购物车创建订单）
    
    Args:
        email: 用户邮箱
        name: 收货人姓名
        street_address: 街道地址
        city: 城市
        zip_code: 邮编
    
    Returns:
        订单创建结果
    """
    if _api_client is None:
        return "Error: API client not initialized"
    
    try:
        shipping_address = {
            "name": name,
            "street_address": street_address,
            "city": city,
            "zip_code": zip_code
        }
        
        result = _api_client.place_order(email, shipping_address)
        data = _extract_data(result)
        
        # API返回可能是 {order: {...}} 或 {order_id: ...}
        if isinstance(data, dict):
            order_id = data.get("order_id") or data.get("order", {}).get("order_id", "")
            if order_id:
                return f"🎉 订单创建成功！\n订单号: {order_id}"
        return "✅ 订单创建成功！"
    except Exception as e:
        return f"❌ 下单出错: {str(e)}"


# 导出所有工具
def get_all_tools():
    """获取所有工具列表"""
    return [
        search_products_tool,
        get_product_detail_tool,
        view_cart_tool,
        add_to_cart_tool,
        remove_from_cart_tool,
        place_order_tool
    ]
