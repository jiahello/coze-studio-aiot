# Coze-Studio-AIOT

[![Coze-Studio-AIOT](https://img.shields.io/badge/Coze-Studio--AIOT-blue)](<[https://github.com/coze-dev/coze-studio](https://github.com/jiahello/Coze-Studio-AIOT)>)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)

**Coze-Studio-AIOT** 是一个基于 Coze Studio 开源框架构建的 AIoT (AI + IoT) 扩展平台，专为智能硬件开发者设计。本项目将 Coze Studio 的强大 AI 工作流能力与物联网设备管理相结合，支持多种硬件平台（包括但不限于 ESP32 系列设备），而"小智"仅作为其中一款设备的参考实现。

## 核心功能

- **可视化 AI 工作流构建**：通过拖拽节点自由编排任意具有 workflow 的 AI Agent，无需编写代码即可创建复杂的 AI 应用逻辑
- **多设备支持框架**：提供统一的设备接入标准，支持 ESP32、树莓派等多种硬件平台，小智只是其中一款参考设备
- **模型服务集成**：预置主流模型服务模板（如火山方舟、OpenAI 等），简化模型配置流程
- **API 与 SDK 部署**：可将 AI 应用部署为 API 或 Web SDK，轻松集成到其他应用程序中
- **知识库整合**：支持将领域知识集成到 AI 工作流中，提升设备交互的智能化水平
- **设备调试平台**：内置设备调试工具和管理界面，涵盖从开发到部署的完整流程

## 与小智 ESP32 服务器的区别

| 特性     | 小智 ESP32 服务器 | Coze-Studio-AIOT               |
| -------- | ----------------- | ------------------------------ |
| 架构     | 单一设备专用      | 多设备通用框架                 |
| 扩展性   | 有限              | 支持插件式扩展新设备类型       |
| AI 能力  | 基础 AI 功能      | 完整 Coze Studio AI 工作流能力 |
| 部署方式 | 设备端独立运行    | 支持云-边-端协同部署           |
| 开发模式 | 代码开发为主      | 可视化低代码/无代码开发        |

## 快速开始

### 前置条件

- Node.js v18+
- Python 3.8+
- 有效的模型服务 API 密钥（如 OpenAI）

### 安装步骤

1. **克隆项目**

   ```bash
   git clone https://github.com/your-username/Coze-Studio-AIOT.git
   cd Coze-Studio-AIOT
   ```

2. **配置模型服务**
   在 `backend/conf/model/template` 目录下选择适合的模型模板文件，根据您的服务提供商进行配置

3. **启动开发环境**

   ```bash
   # 安装依赖
   npm install
   pip install -r requirements.txt

   # 启动服务
   npm run dev
   ```

4. **访问管理界面**
   打开浏览器访问 `http://localhost:3000`，使用可视化界面创建您的第一个 AIoT 工作流

## 设备集成示例：添加小智 ESP32

```markdown
1. 在设备管理界面选择"添加新设备"
2. 选择"ESP32"设备模板
3. 配置设备参数（IP 地址、认证信息等）
4. 通过可视化工作流设计器创建设备专属 AI 能力
5. 部署工作流到设备
```

## 架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                   Coze-Studio-AIOT 平台                       │
├─────────────┬───────────────────────┬───────────────────────┤
│  设备接入层  │     AI 工作流引擎      │     管理控制台       │
│ (支持多硬件) │ (可视化编排/执行环境)  │ (设备管理/调试工具)   │
└─────────────┴───────────────────────┴───────────────────────┘
            ▲                   ▲                   ▲
            │                   │                   │
┌───────────┴───┐     ┌─────────┴─────┐     ┌───────┴───────────┐
│  小智 ESP32   │     │  其他 AIoT 设备  │     │  开发者/用户界面   │
└───────────────┘     └───────────────┘     └───────────────────┘
```

## 贡献指南

欢迎贡献！请遵循以下步骤：

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

## 许可证

本项目基于 Apache License 2.0 开源许可发布 - 详见 [LICENSE](LICENSE) 文件。

## 相关资源

- [Coze Studio 官方文档](https://github.com/coze-dev/coze-studio/wiki)
- [AIoT 设备开发指南](https://github.com/coze-dev/coze-studio/wiki/2.-%E5%BF%AB%E9%80%9F%E5%BC%80%E5%A7%8B)
- [Coze 工作流设计教程](https://www.coze.com/guides/features)
