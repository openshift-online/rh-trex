# TRex Environments, Dependency Injection and Service Locator 

When working with different application environments (dev, prod...) a best practice is to keep **Environment parity** and have all environments as similar as possible. TRex promotes this idea by configuring by default the application with production settings, but allowing each environment to override specific components.

The order of the overrides is for each type of component:
- Config (including Flags)
- Database
- Clients
- Services

If an environment requires a different implementation or configuration  (e.g. using an in-memory DAO for development), then the **Environment framework** is used. Each environment configuration is defined in a separated `e_<environment>.go` file that implements the **Visitor pattern**. Different `Visit<Type>` functions are invoked to create overrides for each type of object.

TRex initializes objects by default using the `Init()` function of go packages. For an `<entity>`, the initialization happens in the `<entity>/plugin.go` with production values and dependencies. For services, the registration uses a `<Entity>ServiceLocator` pattern that returns a function, so the instantiation of service objects is delayed until the service is used by the application code (Lazy instantiation)

This allows services to inject dependencies of other services and the instantiation order will be resolved at runtime.

In the following example the `DinosaurService` is configured with `AdvisoryLock`, `DinosaurDao` and `EventService`

```go
func NewDinosaurServiceLocator(env *environments.Env) DinosaurServiceLocator {
	return func() services.DinosaurService {
		return services.NewDinosaurService(
			db.NewAdvisoryLockFactory(env.Database.SessionFactory),
			dao.NewDinosaurDao(&env.Database.SessionFactory),
			events.EventService(&env.Services),
		)
	}
}
```

Note that at this point, the database session factory has been already initialized by `db/db_session/default.go` `Init` function.
But the `EventService` is being retrieved from the environment, so we could be using one different from the default.

In the TRex application example, the actual instantiation of the `DinosaurService` happens:
- Initializing the controllers (the background process that listens to dinosaur changes)
- Configuring the HTTP routes, which use the handlers that use the services

Note that in this case there are two different instances of the `DinosaurService` type. If we want to have a singleton (only one instance per application), we could do that in the function that resolves the service from the registry.
