"""State definitions and node functions for the product listing agent.

This module defines the state graph for the product listing workflow using LangGraph.
"""

from typing import TypedDict, Annotated, Literal, Optional
from operator import add
from langchain_core.messages import BaseMessage, HumanMessage, AIMessage, SystemMessage


class ProductListingState(TypedDict):
    """State for the product listing workflow.
    
    Attributes:
        messages: List of conversation messages
        user_input: Original user input about the product
        product_data: Current product data being constructed
        validation_status: Whether the product data is validated (pending/approved/rejected)
        error_message: Error message from API if creation failed
        retry_count: Number of retries attempted
        next_step: Next step to execute in the workflow
    """
    messages: Annotated[list[BaseMessage], add]
    user_input: str
    product_data: Optional[dict]
    validation_status: Literal["pending", "approved", "rejected"]
    error_message: Optional[str]
    retry_count: int
    next_step: Literal["complete_info", "confirm", "create", "retry", "end", "waiting_for_user"]


def create_initial_state(user_input: str) -> ProductListingState:
    """Create initial state for the workflow.
    
    Args:
        user_input: User's initial input about the product to create
        
    Returns:
        Initial state dictionary
    """
    return ProductListingState(
        messages=[HumanMessage(content=user_input)],
        user_input=user_input,
        product_data=None,
        validation_status="pending",
        error_message=None,
        retry_count=0,
        next_step="complete_info"
    )


# System prompts for different stages
COMPLETE_INFO_SYSTEM_PROMPT = """你是一个专业的电商商品上架助手。你的任务是将用户的商品描述转换为完整的JSON数据结构。

**核心原则：即使信息不足，也必须生成完整的商品数据！**
- 缺少价格？生成合理的市场价格
- 缺少规格？根据常见配置自动补充
- 缺少描述？生成专业的商品介绍
- 用户后续可以修改任何数据，所以先生成完整结构最重要！

**重要：你必须严格按照以下JSON格式输出，不要有任何其他文字！**

```json
{
  "spu": {
    "brand_id": 整数,
    "category_id": 整数,
    "name": "商品名称",
    "sub_title": "副标题（自动生成吸引人的描述）",
    "main_image": "https://example.com/product/商品名缩写/main.jpg",
    "sort": 0,
    "service_bits": 0
  },
  "skus": [
    {
      "sku_code": "SKU编码",
      "name": "SKU名称",
      "sub_title": "SKU副标题（描述规格特点）",
      "main_image": "https://example.com/product/商品名缩写/sku1.jpg",
      "price": "价格字符串如99.99",
      "market_price": "市场价字符串（比price高20-30%）",
      "stock": 库存整数,
      "sku_spec_data": "规格JSON字符串如{\"color\":\"黑色\",\"size\":\"256GB\"}"
    }
  ],
  "detail": {
    "description": "<p>商品详细描述，包含特点、功能、使用场景等</p>",
    "images": ["https://example.com/product/商品名缩写/detail1.jpg", "https://example.com/product/商品名缩写/detail2.jpg"],
    "videos": [],
    "market_tag_json": "{}",
    "tech_tag_json": "{\"参数名\":\"参数值\"}"
  }
}
```

**可用工具：**
- get_categories_tool: 获取分类列表，找到合适的category_id
- get_brands_tool: 获取品牌列表，找到合适的brand_id

**步骤：**
1. 使用 get_brands_tool 查询品牌（如"小米"、"Apple"等）获取 brand_id
2. 使用 get_categories_tool 查询分类（如"手机"、"家电"等）获取 category_id
3. 为每个规格创建一个SKU对象
4. 生成合理的示例图片URL（使用 https://example.com/product/... 格式）
5. 如果用户没提供价格，设置一个合理的默认值
6. 为market_price设置比price高20-30%的价格
7. SKU编码格式：品牌缩写-商品型号-规格（如MI-HUM-5L）

**关键规则：**
- 价格必须是字符串："99.99" 不是 99.99
- 至少要有1个SKU，多规格就多个SKU
- **图片URL必须生成示例链接，不要留空！**格式：https://example.com/product/商品拼音/图片名.jpg
- 为detail.images至少生成2个示例图片URL
- sub_title要生成吸引人的描述，不要留空
- detail.description要生成HTML格式的详细描述
- tech_tag_json要包含商品的技术参数
- 如果用户没提供库存，默认设为100
- **重要：无论用户信息是否完整，都必须生成完整的JSON数据！缺少的信息用合理的默认值或常见配置填充！**
- 用户可以在下一步修改任何数据，所以现在先把结构生成出来最重要
- 最后只输出JSON，不要有解释文字！
"""

VALIDATION_SYSTEM_PROMPT = """你是一个商品数据验证助手。

用户会提供修改意见或确认信息。

如果用户确认商品信息无误，返回JSON: {"action": "approved", "data": <原始商品数据>}
如果用户提出修改，根据修改意见更新商品数据，返回JSON: {"action": "rejected", "data": <更新后的商品数据>, "reason": "用户要求的修改内容"}
"""

ERROR_RETRY_SYSTEM_PROMPT = """你是一个商品数据修复助手。

商品创建失败了，错误信息如下：
{error_message}

请分析错误原因，修正商品数据中的问题：
1. 检查必需字段是否完整
2. 检查数据类型是否正确（价格是字符串，ID是整数等）
3. 检查字段值是否合法

返回修正后的完整商品数据JSON。
"""


def format_product_data_for_display(data: dict) -> str:
    """Format product data for user-friendly display.
    
    Args:
        data: Product data dictionary
        
    Returns:
        Formatted string representation
    """
    import json
    result = ["📦 商品信息预览：\n"]
    
    if "spu" in data:
        spu = data["spu"]
        result.append("🏷️ SPU信息：")
        result.append(f"  - 名称: {spu.get('name', 'N/A')}")
        result.append(f"  - 副标题: {spu.get('sub_title', 'N/A')}")
        result.append(f"  - 品牌ID: {spu.get('brand_id', 'N/A')}")
        result.append(f"  - 分类ID: {spu.get('category_id', 'N/A')}")
        result.append("")
    
    if "skus" in data:
        result.append(f"📊 SKU信息（共{len(data['skus'])}个）：")
        for i, sku in enumerate(data["skus"], 1):
            result.append(f"  SKU #{i}:")
            result.append(f"    - 编码: {sku.get('sku_code', 'N/A')}")
            result.append(f"    - 名称: {sku.get('name', 'N/A')}")
            result.append(f"    - 价格: ¥{sku.get('price', 'N/A')}")
            result.append(f"    - 库存: {sku.get('stock', 'N/A')}")
        result.append("")
    
    result.append("\n完整JSON数据：")
    result.append(json.dumps(data, ensure_ascii=False, indent=2))
    
    return "\n".join(result)
