"""API client wrapper for pmall API services.

This module provides a simple wrapper around the pmall API endpoints,
handling authentication and providing convenient methods for product listing operations.
"""

import os
import requests
from typing import Optional, Dict, Any, List
from dotenv import load_dotenv

load_dotenv()


class PmallAPIClient:
    """Client for interacting with the pmall API."""
    
    def __init__(self, base_url: str = None, username: str = None, password: str = None):
        """Initialize the API client.
        
        Args:
            base_url: Base URL of the API (default: http://localhost:8080)
            username: Username for authentication (default: from env or 'piao')
            password: Password for authentication (default: from env or '123456')
        """
        self.base_url = base_url or os.getenv("PMALL_API_URL", "http://localhost:8080")
        self.username = username or os.getenv("PMALL_USERNAME", "piao")
        self.password = password or os.getenv("PMALL_PASSWORD", "123456")
        self.token: Optional[str] = None
        self.user_id: Optional[int] = None
        
    def _make_request(
        self, 
        method: str, 
        endpoint: str, 
        data: Optional[Dict] = None,
        params: Optional[Dict] = None,
        require_auth: bool = False
    ) -> Dict[str, Any]:
        """Make HTTP request to the API.
        
        Args:
            method: HTTP method (GET, POST, etc.)
            endpoint: API endpoint (without base URL)
            data: Request body data
            params: Query parameters
            require_auth: Whether authentication is required
            
        Returns:
            Response JSON data
            
        Raises:
            Exception: If request fails
        """
        url = f"{self.base_url}{endpoint}"
        headers = {"Content-Type": "application/json"}
        
        if require_auth and self.token:
            headers["Authorization"] = f"Bearer {self.token}"
        
        try:
            response = requests.request(
                method=method,
                url=url,
                json=data,
                params=params,
                headers=headers,
                timeout=30
            )
            response.raise_for_status()
            result = response.json()
            
            # Handle wrapped response format {code, message, data}
            if isinstance(result, dict) and "data" in result:
                return result["data"]
            return result
        except requests.exceptions.RequestException as e:
            raise Exception(f"API request failed: {str(e)}")
    
    def login(self) -> Dict[str, Any]:
        """Login to the API and store the authentication token.
        
        Returns:
            Login response with user info and token
        """
        response = self._make_request(
            method="POST",
            endpoint="/login",
            data={"username": self.username, "password": self.password}
        )
        
        if "token" in response:
            self.token = response["token"]
            if "user" in response and "id" in response["user"]:
                self.user_id = response["user"]["id"]
        
        return response
    
    def ensure_authenticated(self):
        """Ensure the client is authenticated, login if not."""
        if not self.token:
            self.login()
    
    def get_categories(self, parent_id: Optional[int] = None) -> List[Dict[str, Any]]:
        """Get product categories.
        
        Args:
            parent_id: Parent category ID (None for root categories)
            
        Returns:
            List of categories
        """
        params = {}
        if parent_id is not None:
            params["parent_id"] = parent_id
            
        response = self._make_request(
            method="GET",
            endpoint="/categories",
            params=params
        )
        return response.get("categories", [])
    
    def get_brands(self, page: int = 1, page_size: int = 100) -> Dict[str, Any]:
        """Get product brands.
        
        Args:
            page: Page number
            page_size: Items per page
            
        Returns:
            Brands list and total count
        """
        response = self._make_request(
            method="GET",
            endpoint="/brands",
            params={"page": page, "page_size": page_size}
        )
        return response
    
    def create_product(
        self,
        spu: Dict[str, Any],
        skus: List[Dict[str, Any]],
        detail: Optional[Dict[str, Any]] = None
    ) -> Dict[str, Any]:
        """Create a new product.
        
        Args:
            spu: SPU (Standard Product Unit) data
            skus: List of SKU (Stock Keeping Unit) data
            detail: Product detail information
            
        Returns:
            Creation response with spu_id and message
        """
        self.ensure_authenticated()
        
        data = {
            "spu": spu,
            "skus": skus,
        }
        if detail:
            data["detail"] = detail
            
        return self._make_request(
            method="POST",
            endpoint="/admin/products",
            data=data,
            require_auth=True
        )
