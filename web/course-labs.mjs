const frame = (node, message, code, values, failed = false) => ({ node, message, code, values, failed });
const choice = (id, label) => ({ id, label });
export const courseLabs = {
  '01': { title: '函数的多返回值与错误分支', source: '01-quickstart/03_function.go', nodes: ['调用者', 'divide(a, b)', 'result 与 err'], options: [choice('normal','10 ÷ 2'),choice('zero','10 ÷ 0'),choice('integer','7 ÷ 2')], note: 'Go 的 int 除法保留整数部分。error 是普通返回值，调用方显式判断；这里逐步展示仓库 divide 函数的两条路径。', steps(option) {
    const a = option === 'integer' ? 7 : 10, b = option === 'zero' ? 0 : 2;
    return [frame(0,`调用 divide(${a}, ${b})，准备接收两个返回值。`,`result, err := divide(${a}, ${b})`, {a,b,result:'未返回',err:'未返回'}),
      frame(1,`判断 b == 0：${b === 0}。`,'if b == 0 {', {a,b,'进入错误分支':b === 0}),
      frame(1,b === 0 ? '没有执行除法；返回零值与错误。' : '执行整数除法，返回结果与 nil。', b === 0 ? 'return 0, errors.New("除数不能为0")' : 'return a / b, nil', {result:b === 0 ? 0 : Math.trunc(a/b),err:b === 0 ? '除数不能为0' : 'nil'}),
      frame(2,b === 0 ? 'err 非 nil，调用方处理错误，不能把 result 当成有效商。' : 'err 为 nil，调用方使用 result。','if err != nil { /* 处理错误 */ } else { /* 使用结果 */ }', {result:b === 0 ? 0 : Math.trunc(a/b),err:b === 0 ? '除数不能为0' : 'nil'}, b === 0)];
  } },
  '02': { title:'一次查询如何穿过 API 分层',source:'02-user-api/controller/user_controller.go',nodes:['HTTP 请求','Controller','Service / Repository','数据库','HTTP 响应'],options:[choice('found','用户存在'),choice('missing','用户不存在'),choice('invalid','ID 不是数字')],note:'按 GetByID 的真实控制流模拟。数据库用户是示例数据，不会连接 PostgreSQL；非法 ID 在 Controller 提前返回，不会进入查询。',steps(option){
    const id=option==='invalid'?'abc':option==='missing'?'999':'1';
    const frames=[frame(0,`客户端查询用户 ${id}。`,`GET /api/users/${id}`,{id,'数据库查询次数':0}),frame(1,'先把路径参数解析成无符号整数。','id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)',{id,'解析结果':option==='invalid'?'失败':'成功'})];
    if(option==='invalid')return [...frames,frame(4,'参数错误，立即返回 400；后续层不会执行。','ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})',{'HTTP':400,'数据库查询次数':0},true)];
    return [...frames,frame(2,'Controller 依赖 Service 接口，再由数据访问层查询。','user, err := c.service.GetUserByID(uint(id))',{id,'调用方向':'Controller → Service → Repository'}),frame(3,option==='missing'?'没有对应记录，得到 gorm.ErrRecordNotFound。':'读取示例用户记录。','// 数据库查询完成',{'数据库查询次数':1,'查询结果':option==='missing'?'gorm.ErrRecordNotFound':'User{ID: 1}'}),frame(4,option==='missing'?'errors.Is 匹配记录不存在，映射为 404。':'查询成功，用户数据编码为 JSON。',option==='missing'?'ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})':'ctx.JSON(http.StatusOK, user)',{'HTTP':option==='missing'?404:200},option==='missing')];
  } },
  '03': { title:'JWT 中间件的放行与中止',source:'03-user-api-advanced/middleware/auth.go',nodes:['请求头','格式检查','解析 Token','业务处理'],options:[choice('valid','有效 Bearer Token'),choice('missing','没有 Authorization'),choice('malformed','缺少 Bearer 前缀'),choice('expired','Token 已过期')],note:'演示 AuthMiddleware 的控制流，不生成或校验真实 JWT。只有格式与 Token 解析都通过，才会设置 userID 并 c.Next()。',steps(option){
    const header=option==='missing'?'空':option==='malformed'?'Token demo':'Bearer demo';
    const frames=[frame(0,'读取 Authorization 请求头。','authHeader := c.GetHeader("Authorization")',{'Authorization':header})];
    if(option==='missing')return [...frames,frame(0,'缺少请求头，返回 401 并中止后续处理。','utils.Unauthorized(c, "请先登录"); c.Abort(); return',{'HTTP':401,'业务处理':'未执行'},true)];
    frames.push(frame(1,'按空格拆分，检查前缀必须是 Bearer。','parts := strings.SplitN(authHeader, " ", 2)',{'前缀':option==='malformed'?'Token':'Bearer'}));
    if(option==='malformed')return [...frames,frame(1,'前缀不匹配，返回 401。','utils.Unauthorized(c, "Token格式错误"); c.Abort(); return',{'HTTP':401,'业务处理':'未执行'},true)];
    frames.push(frame(2,'调用 ParseToken，检查签名与有效期。','claims, err := utils.ParseToken(parts[1])',{'解析结果':option==='expired'?'过期错误':'有效 claims'}));
    if(option==='expired')return [...frames,frame(2,'解析失败，返回 401，不能进入业务处理。','utils.Unauthorized(c, "Token无效或已过期"); c.Abort(); return',{'HTTP':401,'业务处理':'未执行'},true)];
    return [...frames,frame(3,'把身份写入 Gin Context，继续执行后续处理。','c.Set("userID", claims.UserID); c.Next()',{'userID':'1（示例）','业务处理':'放行'})];
  } },
  '04': { title:'gRPC 四种通信模式',source:'04-grpc-demo/proto/demo.proto',nodes:['客户端','RPC 数据流','服务端'],options:[choice('unary','Unary · 一来一回'),choice('server','Server streaming · 一来多回'),choice('client','Client streaming · 多来一回'),choice('bidi','Bidirectional · 双向流')],note:'每步是一条示例消息。Proto 的 stream 修饰哪一侧，哪一侧就能发送多条消息；双向流允许独立收发，动画只展示一种可能顺序。',steps(option){
    const flows={unary:[['请求用户 1','→'],['返回用户 1','←']],server:[['请求用户列表','→'],['用户 1','←'],['用户 2','←'],['用户 3','←']],client:[['上传用户 1','→'],['上传用户 2','→'],['客户端关闭发送方向','→'],['汇总 count=2','←']],bidi:[['你好','→'],['你好，客户端','←'],['问题 A','→'],['问题 B','→'],['答复 A','←'],['答复 B','←']]};
    const signatures={unary:'rpc GetUser(GetUserRequest) returns (GetUserResponse);',server:'rpc ListUsers(ListUsersRequest) returns (stream ListUsersResponse);',client:'rpc UploadUsers(stream UploadUsersRequest) returns (UploadUsersResponse);',bidi:'rpc Chat(stream ChatMessage) returns (stream ChatMessage);'};
    let sent=0,received=0;
    return flows[option].map(([message,direction])=>{
      if(direction==='→'&&!message.includes('关闭'))sent++;if(direction==='←')received++;
      return frame(direction==='→'?2:0,`${direction==='→'?'客户端 → 服务端':'服务端 → 客户端'}：${message}`,signatures[option],{'方向':direction,'客户端消息数':sent,'服务端消息数':received,'本次消息':message});
    });
  } },
  '05': { title:'订单服务如何依赖用户服务',source:'05-grpc-microservice/order-service/main.go',nodes:['客户端','OrderService','UserService','订单存储'],options:[choice('create','创建订单 · 用户存在'),choice('missing','创建订单 · 用户不存在'),choice('cancel','创建订单 · 下游因取消失败'),choice('degrade','查询订单 · 用户服务不可用')],note:'对应 CreateOrder 与 GetOrder 的不同处理。当前示例把创建订单时的所有下游错误映射为 NotFound，包括取消错误；动画如实展示这一简化。查询已有订单时则用 unknown 用户名降级。',steps(option){
    const query=option==='degrade';
    const frames=[frame(0,query?'客户端查询一条已存在订单。':'客户端发起创建订单请求。',query?'GetOrder(ctx, req)':'CreateOrder(ctx, req)',{'操作':query?'查询':'创建','订单写入次数':0}),frame(1,query?'从内存读取订单，再查询最新用户名。':'订单服务透传请求 Context，验证用户。',query?'order, ok := s.orders[req.Id]':'s.userClient.GetUser(ctx, &userpb.GetUserRequest{Id: req.UserId})',{'连接':'复用 UserServiceClient','Context':'透传'})];
    frames.push(frame(2,option==='create'?'返回示例用户 Alice。':option==='missing'?'用户不存在，返回错误。':option==='cancel'?'请求 Context 取消，下游返回错误。':'用户服务不可用，返回错误。','userResp, err := s.userClient.GetUser(ctx, ...)',{'下游结果':option==='create'?'User{Name: Alice}':'error'},option!=='create'));
    if(query)return [...frames,frame(1,'用户名降级为 unknown，已有订单仍可返回。','userName := "unknown"',{'订单':'已有记录','user_name':'unknown','订单写入次数':0})];
    if(option!=='create')return [...frames,frame(1,'创建提前返回；不会加锁或写入订单。','return nil, status.Errorf(codes.NotFound, "user not found: %d", req.UserId)',{'返回码':'NotFound','订单写入次数':0},true)];
    return [...frames,frame(3,'加锁，写入订单核心字段并递增 ID；存储不冗余用户名。','s.orders[s.nextID] = order; s.nextID++',{'订单写入次数':1,'存储':'ID / UserId / Product / Quantity'}),frame(0,'构建响应时补充用户名并返回。','UserName: userResp.User.Name',{'返回':'OrderView','user_name':'Alice（示例）'})];
  } },
  '07': { title:'令牌桶如何允许突发、限制速率',source:'07-backend-patterns/rate_limiter/main.go',nodes:['请求','令牌桶','允许 / 拒绝'],options:[choice('burst','同时发出 6 个请求'),choice('paced','每隔 500ms 发一个请求'),choice('fast','每隔 100ms 发一个请求')],note:'模拟 demoAllow 的 rate.NewLimiter(2, 3)：初始 3 个令牌，每秒补 2 个，最多 3 个。时间只随动画步骤推进，不受真实播放速度影响；Allow 没有令牌就立即拒绝。',steps(option){
    const interval=option==='paced'?0.5:option==='fast'?0.1:0;
    let tokens=3,pass=0,reject=0;
    const frames=[frame(1,'初始桶满：有 3 个令牌，可以允许一次突发。','limiter := rate.NewLimiter(2, 3)',{'模拟时间':'0.0s','令牌':3,'通过':0,'拒绝':0})];
    for(let i=0;i<6;i++){
      tokens=Math.min(3,tokens+(i?interval*2:0));const before=tokens;const allowed=tokens>=1;
      if(allowed){tokens--;pass++;}else reject++;
      frames.push(frame(allowed?2:1,`请求 ${i+1}：${allowed?'消耗 1 个令牌，允许通过。':'不足 1 个令牌，立即拒绝。'}`,'allowed := limiter.Allow()',{'模拟时间':`${(i*interval).toFixed(1)}s`,'请求前令牌':Number(before.toFixed(2)),'令牌':Number(tokens.toFixed(2)),'通过':pass,'拒绝':reject},!allowed));
    }
    return frames;
  } },
  '08': { title:'错误包装后，根因还在吗？',source:'08-engineering/01-errors/main.go',nodes:['底层错误','包装层','调用方判断'],options:[choice('wrap','使用 %w 包装'),choice('string','改成 %v 字符串化'),choice('typed','errors.As 提取类型')],note:'%w 保留错误链，%v 只保留错误文本。%v 是对课程代码的对照改造。errors.Is 匹配哨兵，errors.As 提取匹配的错误类型，不依赖字符串判断。',steps(option){
    if(option==='typed')return [frame(0,'validateAge(200) 返回自定义 ValidationError。','return &ValidationError{Field: "age", Msg: "out of range: 200"}',{'类型':'*ValidationError','Field':'age'}),frame(2,'errors.As 找到匹配的错误类型，将值写入 ve。','var ve *ValidationError; errors.As(verr, &ve)',{'匹配':true,'ve.Field':'age','ve.Msg':'out of range: 200'})];
    const wrap=option==='wrap';return [frame(0,'findUser(-1) 触发资源不存在的哨兵错误。','ErrNotFound = errors.New("resource not found")',{'根因':'ErrNotFound'}),frame(1,wrap?'包装附加上下文，同时保留原始错误。':'仅把错误转换成文本，丢失可遍历的错误链。',`fmt.Errorf("findUser(id=%d): %${wrap?'w':'v'}", id, ErrNotFound)`,{'文本':'findUser(id=-1): resource not found','保留错误链':wrap}),frame(2,wrap?'穿透包装，仍然能匹配 ErrNotFound。':'文本虽然一样，但不能匹配原始哨兵错误。','errors.Is(err, ErrNotFound)',{'匹配':wrap,'判断方式':'错误身份，而非字符串'},!wrap)];
  } },
};
export function labForCourse(id) {
  if(id==='06')return {panel:'channel',title:'Channel 的发送与接收'};
  if(id==='09')return {panel:'gateway',title:'SSE 与上游切换'};
  const lab=courseLabs[id];
  if(!lab)throw new Error(`No teaching lab for course ${id}`);
  return {panel:'course',title:lab.title,lab};
}
