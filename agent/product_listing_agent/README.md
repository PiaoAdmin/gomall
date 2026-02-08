# Product Listing Agent

基于 LangGraph 的智能商品上架助手，支持自然语言输入，自动补全商品信息。

## 功能特性

- 🤖 **智能补全**：基于用户简单描述，自动补全完整的商品数据结构
- 🔄 **状态管理**：使用 LangGraph 实现清晰的状态转换流程
- ✅ **人机交互**：支持用户确认和修改商品信息
- 🔧 **自动修复**：API 调用失败时自动分析错误并重试（最多3次）
- 🛠️ **工具调用**：自动查询品牌、分类等辅助信息

## 工作流程

```
用户输入商品描述
    ↓
LLM 补全商品信息（可调用工具查询品牌/分类）
    ↓
展示商品信息给用户确认
    ↓
用户确认/修改
    ↓
    ├─ 确认 → 调用 API 创建商品
    │          ↓
    │       成功 → 结束
    │          ↓
    │       失败 → LLM 分析错误并修复
    │          ↓
    │       重试创建（最多3次）
    │
    └─ 修改 → LLM 更新数据 → 重新确认
```

## 环境要求

- Python >= 3.12
- pmall API 服务运行在 http://localhost:8080

## 安装

确保项目依赖已安装：

```bash
cd /home/piao/workspace/piao_blog/p_mall/agent
pip install -e .
```

## 配置

创建 `.env` 文件或设置环境变量：

```bash
# pmall API 配置
PMALL_API_URL=http://localhost:8080
PMALL_USERNAME=piao
PMALL_PASSWORD=123456

# LLM 配置（使用 OpenAI 兼容 API）
OPENAI_API_KEY=your_api_key_here
OPENAI_API_BASE=https://dashscope.aliyuncs.com/compatible-mode/v1  # 可选
OPENAI_MODEL=qwen-plus  # 默认模型
```

## 使用方法

### 启动交互式界面

```bash
python -m product_listing_agent.main
```

### 使用示例

```
👤 你: 帮我添加一个iPhone 15 Pro Max，256GB，价格8999元，库存100

🤖 Agent: 正在分析商品信息...

🤖 Agent: 📦 商品信息预览：

🏷️ SPU信息：
  - 名称: iPhone 15 Pro Max
  - 副标题: 256GB 钛金属
  - 品牌ID: 1
  - 分类ID: 10

📊 SKU信息（共1个）：
  SKU #1:
    - 编码: APPLE-IP15PM-256
    - 名称: iPhone 15 Pro Max 256GB
    - 价格: ¥8999.00
    - 库存: 100

请确认以上商品信息是否正确？
- 输入「是」或「确认」继续创建
- 输入修改意见进行调整

👤 你: 价格改成8888

🤖 Agent: 已根据您的要求更新：
📦 商品信息预览：
...
  - 价格: ¥8888.00
...

👤 你: 确认

🤖 Agent: ✅ 商品创建成功！
SPU ID: 12345
商品创建成功
```

## 状态定义

### ProductListingState

- `messages`: 对话消息列表
- `user_input`: 用户原始输入
- `product_data`: 当前构建的商品数据
- `validation_status`: 验证状态（pending/approved/rejected）
- `error_message`: API 错误信息
- `retry_count`: 重试次数
- `next_step`: 下一步操作（complete_info/confirm/create/retry/end）

## 工具列表

### get_categories_tool
获取商品分类列表，支持按父分类查询。

### get_brands_tool
获取品牌列表，支持分页。

### create_product_tool
创建商品，需要完整的 SPU、SKU 和详情数据。

## 数据结构

### SPU (Standard Product Unit)

```json
{
  "brand_id": 1,           // 品牌ID（必需）
  "category_id": 10,       // 分类ID（必需）
  "name": "商品名称",       // 名称（必需）
  "sub_title": "副标题",   // 副标题（可选）
  "main_image": "图片URL", // 主图（可选）
  "sort": 0,              // 排序（可选）
  "service_bits": 0       // 服务标识（可选）
}
```

### SKU (Stock Keeping Unit)

```json
{
  "sku_code": "BRAND-MODEL",  // SKU编码（必需）
  "name": "SKU名称",          // 名称（必需）
  "sub_title": "副标题",      // 副标题（可选）
  "main_image": "图片URL",    // 主图（可选）
  "price": "99.99",          // 价格字符串（必需）
  "market_price": "199.99",  // 市场价（可选）
  "stock": 100,              // 库存（必需）
  "sku_spec_data": "{}"      // 规格JSON（可选）
}
```

## 注意事项

1. 价格必须使用字符串格式（如 "99.99"）
2. 至少需要提供 1 个 SKU
3. SKU 编码要唯一，建议包含品牌和型号
4. 创建失败会自动重试最多 3 次
5. 确保 API 服务已启动且可访问

## 目录结构

```
product_listing_agent/
├── __init__.py          # 包初始化
├── api_client.py        # API 客户端封装
├── tools.py            # LangChain 工具定义
├── state.py            # 状态和提示词定义
├── graph.py            # LangGraph 工作流
├── main.py             # 主程序入口
└── README.md           # 本文档
```

## 开发

### 添加新工具

1. 在 `api_client.py` 中添加 API 方法
2. 在 `tools.py` 中创建对应的 `@tool` 函数
3. 将工具添加到 `ALL_TOOLS` 列表
4. 更新相关提示词以引导 LLM 使用新工具

### 修改状态流程

编辑 `graph.py` 中的节点函数和路由逻辑。

### 自定义提示词

修改 `state.py` 中的 `*_SYSTEM_PROMPT` 常量。

## 故障排除

### 无法连接 API
- 检查 API 服务是否启动
- 确认 `PMALL_API_URL` 配置正确

### 登录失败
- 验证用户名密码是否正确
- 检查用户是否已注册

### LLM 调用失败
- 确认 `OPENAI_API_KEY` 已设置
- 检查网络连接和 API 额度

### 创建商品失败
- 查看错误信息，Agent 会尝试自动修复
- 检查品牌ID和分类ID是否存在
- 确认必需字段是否完整

## License

MIT
