# 交互学习实验室

纯静态 HTML/CSS/JavaScript，零前端依赖。入口 `index.html`。

## 本地预览

在仓库根目录执行：

```bash
python3 -m http.server 4174 --bind 127.0.0.1 --directory web
```

访问 <http://127.0.0.1:4174/>。可用 `#course-01` 至 `#course-09` 深链接课程；初次进入选中基础语法课。切换课程只展示对应实验，并停止、重置上一课播放。

## 功能与边界

- 九阶段课程：Java 对照、阶段目标、真实源码/练习入口、可复制的本地命令。
- 第 1–5、7、8 课分别演示函数返回值、API 分层、JWT 鉴权、gRPC 四种通信模式、微服务调用、令牌桶与错误链；支持场景选择、播放、前后单步和重置。
- Channel 实验：容量 0/2/4、发送阻塞、接收阻塞、FIFO、关闭后排空与零值。每侧只模拟一个 goroutine；存在阻塞发送时禁止关闭，解释对应 panic 风险。
- 网关实验：正常流、503 fallback、400 停止、输出后读错误、客户端取消。支持播放/暂停、单步和重置；切换实验或隐藏页面会暂停。
- 动画由 `models.mjs` 与 `course-labs.mjs` 中的确定性模型驱动，不运行 Go、不发起上游网络请求，不展示伪造的性能数据。真实流式实现进入读流后的错误不会 fallback，正常 EOF 的行为不属于“读错误”场景。
- 学习进度使用 localStorage，仅保存在当前浏览器，以用户自评为准；存储不可用时仍可学习。
- 键盘可操作全部控件，尊重 reduced-motion 设置。

## 维护与验证

`courses.mjs` 的路径与课程目标必须对照根 README 和源码。模块源码仍是事实来源。新增/修改场景需同步模型、说明和测试。

```bash
node --test web/models.test.mjs
node --check web/app.mjs
```

测试还会验证所有课程目录、练习和答案链接在本地存在。

## 静态托管

将 `web/` 的内容作为站点根目录即可，资源均使用相对路径，也支持项目子路径。

GitHub Pages 由 `.github/workflows/pages.yml` 部署：推送 main 后先运行模型与课程链接测试，再上传站点资源到 Pages。PR 只验证，不部署。

在线地址：<https://hwj123hwj.github.io/go-learning-for-java-dev/>。
