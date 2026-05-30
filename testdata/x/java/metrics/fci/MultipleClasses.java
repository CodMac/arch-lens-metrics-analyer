package metrics.fci;

class ClassA {
    public void mA() {
        // CC = 1
    }
}

public class MultipleClasses {
    public void mB() {
        // CC = 1 + 1 (if) = 2
        if (true) {
            System.out.println("B");
        }
    }
}

class ClassC {
    public void mC() {
        // CC = 1 + 1 (while) = 2
        while(true) {
            break;
        }
    }
}
// Expected FCI = 1 + 2 + 2 = 5
