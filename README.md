# bili
## 服务划分：

1. **API Gateway**
   - 负责外部 HTTP 请求路由、统一鉴权（也可委托 Auth 微服务 验证 JWT）、请求限流、流量管控、静态资源代理、TLS 终端。
   - 对外暴露 `/api/videos/*`、`/api/users/*` 等路径（将前端现有路径映射到对应微服务）。
2. **Auth Service（认证/用户会话）**
   - 注册 / 登录 / Token 签发 & 验证 。
   - 暴露 HTTP REST（供网关/前端调用） + gRPC/HTTP 内部接口。
3. **User Service（用户资料 / 关注 /资料展示）**
   - 管理用户公开信息、头像上传（调用 OSS via Media Service）等。
4. **Media Service（媒体存储/处理）**
   - 负责视频文件分片接收（可选直接签名直传）、OSS Multipart Upload、合并、转码触发（可接入转码队列/FFmpeg）、封面截图、存储 URL。
   - 与 OSS（阿里云 OSS）紧密耦合；对外提供 presigned URL、合并回调、文件状态查询。
5. **Video Service（视频业务）**
   - 负责视频资源的业务信息（title/description/owner/统计）、CRUD、播放统计（接入 Analytics Service）。
   - 通过事件总线订阅 Media 的“转码完成/资源就绪”事件以把视频状态改为 ready。
6. **Comment Service（评论/回复/点赞）**
   - 评论 CRUD、回复树、评论点赞/踩、分页、查询。
7. **Interaction Service（点赞 / 收藏 / 收藏夹）**
   - 可把点赞/收藏独立为服务，方便横向扩展并将频繁变化数据与主视频表分离。
8. **Notification / Message Service（通知）**
   - 关注、评论回复通知，异步发送邮件/站内信/push。
   - 非阻塞，设计为消费者（Kafka/RabbitMQ）。
9. **Analytics Service（播放统计 / PV / UV）**
   - 高吞吐事件收集（play events），通过缓存(Redis)缓冲，周期性写入 ClickHouse 或 MySQL 聚合表。
   - 非阻塞，设计为消费者（Kafka/RabbitMQ）。
10. **Admin Service （管理与计划任务）**
    - 媒体清理、数据回填、统计汇总。