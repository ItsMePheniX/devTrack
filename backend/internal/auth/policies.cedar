//policy: Any user with the 'dev' role is permitted to perform all actions
permit (
    principal,
    action,
    resource
)
when { 
    principal has role && principal.role=="dev"
};