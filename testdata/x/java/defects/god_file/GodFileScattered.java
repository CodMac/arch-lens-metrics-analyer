package com.test.defects.god_file;

import com.order.Order;
import com.user.User;
import com.inventory.Stock;
import com.finance.Invoice;
import com.shipping.Delivery;

/**
 * 场景 2: 逻辑散乱文件 (GodFileScattered)
 * 预期：NDE > 15 (定义了 16 个以上的类/接口), CDC > 4 (跨越了 5 个不同的业务包)
 */
public class GodFileScattered {
    class Entity1 {}
    class Entity2 {}
    class Entity3 {}
    class Entity4 {}
    class Entity5 {}
    class Entity6 {}
    class Entity7 {}
    class Entity8 {}
    class Entity9 {}
    class Entity10 {}
    class Entity11 {}
    class Entity12 {}
    class Entity13 {}
    class Entity14 {}
    class Entity15 {}
    class Entity16 {}
}

interface Service1 {}
interface Service2 {}
interface Service3 {}
