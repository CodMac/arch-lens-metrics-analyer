package com.example.metrics;

// --- 社区 A (Order) ---
class OrderService { void process() {} }
class OrderEntity {}

// --- 社区 B (Inventory) ---
class InventoryManager { void deduct() {} }

// --- 社区 C (User) ---
class UserManager { void getUser() {} }

// --- 社区 D (Finance) ---
class InvoiceService { void issue() {} }

// --- 测试目标：上帝类 ---
public class CrossDomainProcessor {
    // 域内依赖 (不计入 CDC)
    private OrderService order;

    // 跨域依赖 (计入 CDC)
    private InventoryManager inv;
    private UserManager user;
    private InvoiceService finance;

    public void executeAll() {
        order.process(); // A -> A
        inv.deduct();    // A -> B (CDC +1)
        user.getUser();  // A -> C (CDC +1)
        finance.issue(); // A -> D (CDC +1)
    }
}