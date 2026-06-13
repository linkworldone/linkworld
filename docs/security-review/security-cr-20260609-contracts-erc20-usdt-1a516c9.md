# 代码安全审查报告（内部模板）

**审查对象**：LinkWorld 合约层 ERC20 USDT 重构（Deposit / Payment / Oracle / TrafficCardNFT / ServiceManager / FeeManager + interfaces + mocks）
**语言 / 框架**：Solidity ^0.8.27 + OpenZeppelin Upgradeable (UUPS) + ReentrancyGuardTransient + SafeERC20
**审查范围**：`git diff 7ef9677..1a516c9 -- packages/contracts/contracts/`（10 文件，+256/-95 行）
**审查人**：security-reviewer (AI) — 基于内部《安全编码规范》《加密算法与密钥长度推荐》《敏感数据脱敏规则》
**审查日期**：2026-06-09
**业务域**：资金 / 支付结算 / 保证金（链上直分 USDT）
**场景等级**：高（资损敏感）

---

## 一、结论（Conclusion）

**审查结论**：✅ 通过（可放行进 ship）

**一句话摘要**：设计 §6.1 资损 checklist 与 arch-review §七 7 大阻塞项 + A1/A3 在实现代码中逐条落地，未发现新增资损/越权阻塞。

**统计**：

| 严重级 | 数量 | 必须处理时限 |
|--------|------|-------------|
| Critical（致命） | 0 | 立即修复，Block |
| High（高） | 0 | 本次合并前修复，Block |
| Medium（中） | 0 | 本迭代内修复 |
| Low（低） | 1 | 排入后续迭代 |
| Info（提示） | 2 | 知晓即可 |

**误报审查活动（FP Review Activity）**

| 指标 | 数量 |
|---|---|
| 扫描 + 清单命中的候选问题 | 6 |
| 通过 FP Review → Confirmed | 1（Low） |
| 通过 FP Review → Needs Confirmation | 0 |
| 通过 FP Review → Withdrawn（见 §6.2） | 5 |

---

## 二、问题清单（Findings Overview）

| # | 漏洞 | 严重级 | 状态 | CWE | 位置 | OWASP |
|---|------|-------|-----|-----|------|-------|
| SEC-001 | MockUSDT / 非标桩 public `mint` 无访问控制 | Low | Confirmed | CWE-862 | `mocks/MockUSDT.sol:18` | A01 |

---

## 三、详细发现（Detailed Findings）

### [SEC-001] 测试代币 public mint 无访问控制

- **严重级**：Low
- **状态**：Confirmed
- **CWE**：CWE-862 — Missing Authorization
- **OWASP 2021**：A01:2021 Broken Access Control
- **位置**：`packages/contracts/contracts/mocks/MockUSDT.sol:18`、`mocks/NonStandardUSDT.sol:18/29/61`
- **规范条款**：违反《安全编码规范》最小权限原则（任意地址可铸币）

**漏洞描述**

`MockUSDT.mint(to, amount)` 为 external 无修饰符，任意地址可无限铸币。若误将该测试代币当作生产 USDT 部署（deploy.ts 当前在 421614 测试网即用 MockUSDT 作资金通道），任何人可铸币后向 Deposit/Payment 注资或操纵分账。本轮定位为测试/测试网代币，正式上线须替换真实 USDT 地址（deploy.ts:186 注释已声明），故为 Low。

**漏洞代码**

```solidity
/// @notice 公开 mint，测试任意发币
function mint(address to, uint256 amount) external {   // ❌ 无 onlyOwner
    _mint(to, amount);
}
```

**数据流追踪（必填）**

- **Source**：外部调用方任意 EOA（`to`/`amount` 完全可控）
- **中间变换**：无
- **Sink**：`ERC20._mint` 直接增发余额
- **跨越的信任边界**：external 调用 → 余额账本，无权限闸门

**误报分析（False-Positive Analysis，必填）**

1. Source 污点：YES — 调用方与参数完全外部可控
2. Sink 可利用性：YES — 直接无限增发
3. 控制覆盖：NONE — 无 onlyOwner/onlyMinter
4. 环境范围：测试/测试网代币（`mocks/` 目录，design §7.2 明确"本轮无正式 USDT，本地/测试网全用 mock"）→ 因此降级 Low 而非 High
5. 纵深防御：正式上线靠"替换真实 USDT 地址"流程兜底（deploy.ts:186 注释）
6. 规范条款：最小权限原则；mock 故意开放属设计取舍，但需防误用

**可利用 PoC（最小可复现）**

```solidity
// 任意账户：MockUSDT(addr).mint(attacker, 1e30);  // 无限铸币
```

**修复建议（必改代码）**

mock 保持现状可接受（测试便利），但须在生产部署门禁中硬性禁止该合约进入主网 deploy 清单。若希望收敛 mock 风险面，可加 onlyOwner：

```solidity
function mint(address to, uint256 amount) external onlyOwner { _mint(to, amount); }
```

**修复说明**：本轮属设计已接受的测试取舍（mock 故意 public mint 以便补测）。不构成本轮 ❌；仅提示上线 checklist 必须把 MockUSDT/NonStandardUSDT 排除出主网制品。

**验证方法**
- 静态：主网 deploy 脚本不得 import `mocks/`
- 上线门禁：USDT 地址须指向已审计的真实 Tether 合约

**参考**
- OWASP: https://owasp.org/Top10/A01_2021-Broken_Access_Control/
- CWE: https://cwe.mitre.org/data/definitions/862.html

---

## 四、修复优先级（Remediation Priority）

**Top 3 必须立即修复**：无 Critical/High，无阻塞项。

1. **SEC-001** — 上线门禁排除 mock 代币（非代码阻塞，流程项）
2. （无）
3. （无）

**建议修复顺序**：仅 Low/Info，排入上线 checklist 即可。

---

## 五、上线 Gate Checklist

| # | 检查项 | 结论 | 证据 / 备注 |
|---|--------|------|-----|
| 1 | 无硬编码密钥 / 密码 / Secret / Token / 助记词 | 本次无涉 | 合约层无私钥/secret |
| 2 | 所有 SQL 走参数化 | 本次无涉 | 无 SQL |
| 3 | 所有 Controller / Handler 有认证与鉴权 | 通过 | createBill/issueMonthlyTrafficCards onlyOracle；applyTrafficCardToBill onlyOracle；mintTrafficCard onlyOwner；setOperatorPaymentAddress onlyOwner（Payment.sol:74,154；Deposit.sol:86,135；ServiceManager.sol:203） |
| 4 | ID 查询/修改走 owner 校验防水平越权 | 通过 | payBill require `bill.user == msg.sender`（Payment.sol:105）；withdraw 仅操作 msg.sender 自身映射（Deposit.sol:71-79） |
| 5 | 敏感数据日志已脱敏 | 本次无涉 | 事件仅含地址/金额/tokenId，无 PII |
| 6 | CVV / PIN / 私钥 / Seed 不持久化 + 不打印 | 本次无涉 | 无此类数据 |
| 7 | 对称加密 AES-256-GCM | 本次无涉 | 无加密 |
| 8 | 非对称加密 RSA≥2048 / ECC≥256 | 本次无涉 | 无 |
| 9 | 摘要 SHA-256 及以上 | 本次无涉 | 无 |
| 10 | 密码哈希 BCrypt/Argon2id | 本次无涉 | 无 |
| 11 | JWT Secret≥32B 且不硬编码；禁 alg=none | 本次无涉 | 无 JWT |
| 12 | 私钥 / 签名逻辑无外连 | 本次无涉 | 合约无签名机/外连 |
| 13 | TLS 1.2+ 证书验证 | 本次无涉 | 无网络层 |
| 14 | 依赖无已知高危 CVE | 通过 | 仅 OZ contracts/contracts-upgradeable，标准 SafeERC20/ReentrancyGuardTransient |
| 15 | 反序列化白名单 | 本次无涉 | 无 |
| 16 | XML 解析禁 DTD | 本次无涉 | 无 |
| 17 | 文件操作无路径穿越 | 本次无涉 | 无 |
| 18 | 文件上传白名单 | 本次无涉 | 无 |
| 19 | SSRF 防护 | 本次无涉 | 无外部 HTTP |
| 20 | CORS 无 `*`+credentials | 本次无涉 | 无 |
| 21 | 命令执行白名单 | 本次无涉 | 无 |
| 22 | SpEL / 模板不取外部输入 | 本次无涉 | 无 |
| 23 | 异常不泄露栈 | 本次无涉 | revert reason 为静态字符串 |
| 24 | 无裸 println/printStackTrace | 本次无涉 | — |
| 25 | 敏感操作有审计日志 | 通过 | DepositMade/Withdrawn/BillCreated/BillPaid/TrafficCardMinted/TrafficCardApplied/OperatorPaymentAddressSet 事件齐全 |
| 26 | 防重放：幂等 key / nonce / 唯一索引 | 通过 | payBill 靠 `bill.isPaid` 幂等；发卡靠 `getUserCardCount==0` 幂等（Deposit.sol:95） |
| 27 | 金额范围校验防溢出 | 通过 | 0.8.x 默认 checked 算术；calculateFee=(amount*rate)/denom 无溢出；MIN_DEPOSIT 按 decimals 动态算（Deposit.sol:56） |
| 28 | 分布式锁 | 本次无涉 | — |
| 29 | 事务 | 本次无涉 | — |
| 30 | 并发 map 加锁 | 通过 | nonReentrant 覆盖发卡/支付/withdraw 路径 |
| 31 | 依赖固定到唯一不可变版本 + 锁文件 | 本次无涉 | 本次 diff 未改 package.json / 依赖清单文件 |

**Gate 结论**：全部为 `通过` 或 `本次无涉` → **通过 Gate**。

---

## 六、附录

### 6.1 未审查 / 有限审查区域

- `scripts/deploy.ts` 已交叉核对 wiring（setOracle/setPayment 双向、payment.initialize 注入 serviceManager、transferOwnership(deposit)、operator paymentAddress 注入循环、MockUSDT 部署），但属部署脚本非合约代码，未做逐行审计。421614 上链私钥/RPC key 配置为设计已知待配项，不计入本轮 ❌。
- FeeManager / UserRegistry 本次 diff 未改动，仅核对被调用接口（calculateFee / isRegistered）契约一致。

### 6.2 Withdrawn 候选 — 误报披露（False-Positive Disclosure）

| # | 命中的模式 | 位置 | 驳回原因 | 缓解控制（file:line） |
|---|---|---|---|---|
| FP-1 | 发卡路径 external 调用回调重入（`_safeMint` 的 `onERC721Received`） | `Deposit.sol:101-107,134-144` | 发卡路径已加 `nonReentrant`（issueMonthlyTrafficCards / mintTrafficCard），且 NFT.mint 内 `_userCardCount++`(L47) 在 `_safeMint`(L38) 回调后但本合约 guard 已锁；A1 闭合 | `Deposit.sol:86,134`（nonReentrant）+ `TrafficCardNFT.sol:33-50` |
| FP-2 | A3 残留 `_reentrancyGuardInit` 写持久 storage 与 transient guard 不一致 | 原 `Payment.sol:33-38` | 该 init 函数及 initialize 内调用已被删除（diff 确认移除 `_reentrancyGuardInit` + initialize L28 调用），transient guard 由 OZ ReentrancyGuardTransient 原生 tload/tstore 生效，无残留无效代码 | `Payment.sol:46-51`（initialize 已无该调用）|
| FP-3 | payBill 重入（外部 safeTransferFrom 后再操作） | `Payment.sol:102-124` | CEI：`bill.isPaid=true`(L111) 早于两段 safeTransferFrom；且函数 `nonReentrant`(L102)；合约不暂存资金（user→operator / user→platform 直分） | `Payment.sol:102,111` |
| FP-4 | 发卡额度精度 bug `_deposits/100000` 残留 | `Deposit.sol:_mintFor` | 已删除，`dataAmount = trafficCardQuota` 固定额度（v2-C），与存款额/精度解耦；全文无 `/100000` 命中 | `Deposit.sol:101-107` |
| FP-5 | 非标 USDT（无返回值 / 返回 false）导致静默失败资损 | `mocks/NonStandardUSDT.sol` | 资金转移点全部走 `SafeERC20.safeTransfer/safeTransferFrom`，对无返回值代币正常入账、对返回 false 代币整笔 revert；无裸 transfer（grep 验证 0 命中）；fee-on-transfer 靠设计文档约束（design §6.1 B7），属已接受取舍 | `Deposit.sol:58,79`；`Payment.sol:116,120` |

### 6.3 二次审查建议

- 上线门禁须将 `mocks/MockUSDT.sol`、`mocks/NonStandardUSDT.sol` 排除出主网制品，USDT 地址指向真实代币（SEC-001）。
- 若后续 Round 实装 `applyTrafficCardToBill` 真实抵扣语义（本轮为受限桩，不转资金），须对其重新做资损审计。
- fee-on-transfer 代币禁用目前靠文档约束（design §6.1 B7），无链上"实收差额"补偿逻辑——属设计已接受取舍，若未来放开非标代币需补审。

### 6.4 引用规范

- 《安全编码规范》（OWASP Top 10 2021 + CIS Benchmark）
- 《Golang 代码安全审计指南》（本次无涉）
- 《敏感数据脱敏规则》
- 《加密算法与密钥长度使用推荐》
- 《安全评审规范与要求》- 安全红线要求
