"""Pmall API Client for Auto Order Agent

封装商品、购物车、订单相关API接口
"""

import requests
from typing import Optional, Dict, List, Any
import json


class PmallOrderAPIClient:
    """Pmall订单API客户端"""
    
    def __init__(self, base_url: str = "http://localhost:8080"):
        self.base_url = base_url
        self.session = requests.Session()
        self.token: Optional[str] = None
    
    def login(self, username: str, password: str) -> Dict[str, Any]:
        """用户登录"""
        url = f"{self.base_url}/login"
        response = self.session.post(url, json={
            "username": username,
            "password": password
        })
        response.raise_for_status()
        result = response.json()
        
        # API返回格式: {code: 20000, message: "success", data: {token: ..., user: ...}}
        if "data" in result and "token" in result["data"]:
            self.token = result["data"]["token"]
            # API使用cookie 'jwt' 认证，不是Authorization header
            self.session.cookies.set("jwt", self.token)
        elif "token" in result:
            self.token = result["token"]
            self.session.cookies.set("jwt", self.token)
        
        return result
    
    # ==================== 商品相关API ====================
    
    def search_products(
        self, 
        keyword: Optional[str] = None,
        category_id: Optional[int] = None,
        brand_id: Optional[int] = None,
        min_price: Optional[float] = None,
        max_price: Optional[float] = None,
        sort_by: str = "default",
        page: int = 1,
        page_size: int = 10
    ) -> Dict[str, Any]:
        """搜索商品（返回SKU列表）
        
        Args:
            keyword: 搜索关键词
            category_id: 分类ID
            brand_id: 品牌ID
            min_price: 最低价格
            max_price: 最高价格
            sort_by: 排序方式 (default/price_asc/price_desc/sale)
            page: 页码
            page_size: 每页数量
        """
        url = f"{self.base_url}/products/search"
        params = {
            "page": page,
            "page_size": page_size,
            "sort_by": sort_by
        }
        
        if keyword:
            params["keyword"] = keyword
        if category_id:
            params["category_id"] = category_id
        if brand_id:
            params["brand_id"] = brand_id
        if min_price is not None:
            params["min_price"] = str(min_price)
        if max_price is not None:
            params["max_price"] = str(max_price)
        
        response = self.session.get(url, params=params)
        response.raise_for_status()
        return response.json()
    
    def get_product_detail(self, spu_id: int) -> Dict[str, Any]:
        """获取商品详情（SPU + 所有SKU）"""
        url = f"{self.base_url}/products/{spu_id}"
        response = self.session.get(url)
        response.raise_for_status()
        return response.json()
    
    # ==================== 购物车相关API ====================
    
    def get_cart(self) -> Dict[str, Any]:
        """获取购物车详情"""
        url = f"{self.base_url}/cart"
        response = self.session.get(url)
        response.raise_for_status()
        return response.json()
    
    def add_to_cart(self, sku_id: int, quantity: int = 1) -> Dict[str, Any]:
        """添加商品到购物车
        
        Args:
            sku_id: SKU ID
            quantity: 数量
        """
        url = f"{self.base_url}/cart/add"
        response = self.session.post(url, json={
            "sku_id": sku_id,
            "quantity": quantity
        })
        response.raise_for_status()
        return response.json()
    
    def remove_from_cart(self, sku_ids: List[int]) -> Dict[str, Any]:
        """从购物车移除商品"""
        url = f"{self.base_url}/cart/remove"
        response = self.session.post(url, json={
            "sku_ids": sku_ids
        })
        response.raise_for_status()
        return response.json()
    
    def clear_cart(self) -> Dict[str, Any]:
        """清空购物车"""
        url = f"{self.base_url}/cart/clear"
        response = self.session.post(url, json={})
        response.raise_for_status()
        return response.json()
    
    # ==================== 订单相关API ====================
    
    def place_order(
        self, 
        email: str,
        shipping_address: Dict[str, Any]
    ) -> Dict[str, Any]:
        """下单（从购物车创建订单）
        
        Args:
            email: 用户邮箱
            shipping_address: 收货地址 {name, street_address, city, zip_code}
        """
        url = f"{self.base_url}/orders"
        response = self.session.post(url, json={
            "email": email,
            "shipping_address": shipping_address
        })
        response.raise_for_status()
        return response.json()
    
    def list_orders(self) -> Dict[str, Any]:
        """获取用户订单列表"""
        url = f"{self.base_url}/orders"
        response = self.session.get(url)
        response.raise_for_status()
        return response.json()
    
    def cancel_order(self, order_id: str) -> Dict[str, Any]:
        """取消订单"""
        url = f"{self.base_url}/orders/{order_id}/cancel"
        response = self.session.post(url, json={})
        response.raise_for_status()
        return response.json()
