package metrics.fci;

public abstract class Exclusions {
    private String name;

    // CC = 0 (Abstract)
    public abstract void abstractMethod();

    // CC = 0 (Interface method in Java)
    public interface MyInterface {
        void doSomething();
    }

    // CC = 0 (Simple Getter)
    public String getName() {
        return name;
    }

    // CC = 0 (Simple Setter)
    public void setName(String name) {
        this.name = name;
    }

    // CC = 1 (Not a simple getter/setter)
    public String getFormattedName() {
        return "Name: " + name;
    }
}
// Expected FCI = 1
