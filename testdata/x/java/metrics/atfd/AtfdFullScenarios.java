package com.test.atfd;

/* --- 辅助类：外部数据所有者 --- */
class ForeignData {
    public String publicField;
    public static final int CONSTANT = 100;
    private int age;
    private boolean active;

    public int getAge() { return age; }
    public boolean isActive() { return active; }
    public void processData() { /* 业务方法，非 Getter */ }
}

class AnotherForeign {
    public static double price;
    public static double getPrice() { return price; }
}

class BaseClass {
    protected String baseField;
}

/* --- 主类：ATFD 目标计算类 --- */
public class AtfdTarget extends BaseClass {
    private String internalField;
    private ForeignData foreign = new ForeignData();
    private AnotherForeign another = new AnotherForeign();

    /**
     * 场景 1: 自身访问 (Expected ATFD contribution: 0)
     */
    public void selfAccess() {
        this.internalField = "test";
        getSelf();
    }
    private String getSelf() { return internalField; }

    /**
     * 场景 2: 访问 ForeignData (Expected ATFD contribution: 1)
     * 涵盖：字段、get、is、静态、去重
     */
    public void accessForeign() {
        // 1. 直接字段访问
        String f = foreign.publicField;
        
        // 2. Getter 访问
        int a = foreign.getAge();
        
        // 3. Is-Getter 访问
        boolean active = foreign.isActive();
        
        // 4. 静态字段访问
        int c = ForeignData.CONSTANT;
        
        // 5. 非 Getter 调用 (不应计入 ATFD)
        foreign.processData();
        
        // 以上所有操作针对同一个类 ForeignData，ATFD 计数应为 1
    }

    /**
     * 场景 3: 访问 AnotherForeign (Expected ATFD contribution: 1)
     */
    public void accessAnother() {
        double p = AnotherForeign.price;
        double p2 = AnotherForeign.getPrice();
    }

    /**
     * 场景 4: 访问父类数据 (Expected ATFD contribution: 1)
     * 根据协议“非自身类”判定，父类字段属于外部类
     */
    public void accessBase() {
        System.out.println(this.baseField);
    }

    /* 
       总结：
       AtfdTarget 类的总 ATFD 预期值应为 3:
       1. com.test.atfd.ForeignData
       2. com.test.atfd.AnotherForeign
       3. com.test.atfd.BaseClass
    */
}
