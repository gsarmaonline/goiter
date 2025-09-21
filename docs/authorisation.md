# Authorisation

In this section, we will cover how Authorisation works in Goiter.
Every authorisation service has to deal with the following elements:

- Accessor/Actor/User
- Resource/Object
- Action

The underlying statement for an authorisation service is if an accessor should
be allowed to perform an action on the resource.

There are hierarchical concepts which also apply to all the elements in the system.
However, we will try to define a flat structure for now and talk about hierarchical
elements or groups in the future.

## Flat map representation

The easiest way to do this is to have a flat map of all accessors, objects and actions.
To store the mapping, the `RoleAccess` model will be used.

So if we have a mapping with the following

```bash
accessor_id,object_type,object_id,action_type
```

then we can define all possible rules with this structure.

However, the number of rows in the `RoleAccess` model would be tremendously high in this case
and since all the columns are supposed to be indexable, any kind of scan would result in high
resource and time consumption.

For example, if there are 1000 users trying to access 1000 objects, the number of rows would be
a million rows. 1000 users is not that big of a number and anything with more numbers would be
disastrous. On top of it, this model would be used everytime an user tries to access an object.
This means any kind of bottleneck on this model would affect all the APIs.

## Hierarchical representation

To mimic real life scenarios and also to prevent the bloat of the number of rows as mentioned in the
`Flat map` representation, the concept of hierarchies can be brought in.

This representation signifies that every element can be present as a hierarchical `Group` entity.
Any `RoleAccess` rule which matches a group that the accessor belongs to means that they are eligible
to access the object. A group can also belong to another group, which effectively allows it to form a
tree of rules.

Each group can have multiple parents. Each group inherits the properties/rules of the parent groups.
For example, if you don't find the exact match for a specific group, then you can search in the `RoleAccess`
model for the parents of the specific group recursively.

An important assumption is that the depth of recursion to unravel to a matching group is not more than 10.

How does this help? Let's look at an example.

### Example 1

Let's take an example of a case where you want only the finance team to be able to access the billing
section of your app. If your team has 20 people and the number of resources you want to control is
more than 50, then the overall number of rules based on the flat map representation would be 1000.

Now let's look at the example using the grouping or hierarchical representation.

- Create an Object group called `billing` where all the billing objects are placed inside it.
- Create an User group called `finance` where all the finance team members are allocated to.
- Create one rule which allows the `finance` user group to access `billing` objects.

Let's put another restriction where only the executives in the finance team can change the records.
Everybody else in the finance team can only read the records.

- Create an User group called `finance_execs` with `finance` as the parent group.
- Alter the previous rule to allow the `finance` user group to only be able to read the `billing` objects.
- Create one rule which allows the `finance_execs` to perform all operations on the `billing` objects.

### Cons

This system is a whitelisting system. This means that if you don't have any rule which mentions that you
can access the resource, then you can't access the resource.
In future, if we have to support black listing as well, then there may be conflicts between whitelisted
and blacklisted groups and may have to bring on the concept of priorities to the rules.

Another problem is that since an element can belong to groups, for every API, we have to fetch the groups
associated with the elements in a recursive manner till we reach the root or a matching rule. This can lead
to multiple calls, but the scale required for groups would be far lesser than the scale required for the
flat map representation.

## Comparing per element groups or a single group model

In this section, we compare how having a single group model for all element differs from having different
group models for every element.

The first type is having one single model to store Groups of all different elements.
The second type is to have different Group models for every element.

Having the Group elements in different models means that there may be drastically different number of groups
for every element. However, having the right indexes in a single model will result in a similar experience.

Having different Group elements can allow us to store different metadata per element.

## Ownership

The Ownership of an object can be defined as an user which can perform all operations on the object.
How is this different then the above mentioned rules?

In most apps, there are implicit rules that regard the creator of an object to be the owner of the object
and should have all possible actions on the object.

Based on the above representation, if we can define specific `RoleAccess` rules for the user who created the
object. However, this means that we need to have 1 additional rule for every object.

This brings us to the topic of having explicit rules vs implicit rules.
Since this is an app which tries to leverage convention over configuration, implicit rules should be available
to the users as well.

### Implicit rules

Implicit rules should be checked in the `canAccess` method which is the entrypoint for the authorisation checks.
Implicit rules can also be defined as more of a configuration than a rule.

List of implicit rules

- Enable ownership access

## Scope

In every app, there are different kinds of scopes like projects, accounts, etc which provides the encapsulation required
for that level.
For example, if we take `Account` into instance, it's similar to a tenant and the border of this area shouldn't be
crossed by the user accounts.
Then we also have more subscoping in an account using projects. The access rules of a project can be different.
Taking the example of the billing table again, a developer may be able to query some tables but they shouldn't be able
to query the `billing` model of the project.

There are some rules which apply to the entire system as a whole. For example, if I want to create a super user which
can mimic the login of any of the accounts, then there should be a a `RoleAccess` rule applicable for all accounts.

Let's add `scope_type` and `scope_id` to the `RoleAccess` model.
To define rules at the `Account` level for a specific, we add the `scope` as `Account` and `scope_id` of the account.

Now, let's try to define rules on who can create a `Project` and who can add members to the project.
The `canAccess` method receives the user and the object.
From the object, the `scope` of the object can be fetched. For `Project`, the scope has to be defined as `Account` and the
project should also return the account it is linked to via the `scope_id`.
The `canAccess` method searches for the matching rule with `scope_type` of `Account` and the account's `scope_id` as an additional filter.

If we are adding `scope` to the `RoleAccess` model, we also need to add it to the `Group` model. This allows us to
create sub groups for different scopes as well.

Root scopes are the scopes which applies to all objects, irrespective of the tenancy.

### Identifying the Scope of an object

This section explains the process that can be used in identifying the scope of an object.

There are 2 possible options here:

- Every object refers to the actual scope directly
- They refer to the parent object as the scope and the actual scope is recursively
  discovered by going up the stack of the ancestors

## Fetching a list of objects

How do you fetch a list of objects specific to the user in a multi-tenant system?
Most apps add an `owner_id` or `account_id` in every resource and then add it in the `WHERE` clause whenever the list of
objects have to be fetched.

It works for most cases, but few cons are:

- The underlying call to fetch the resources isn't aware of the tenancy. If you forget to omit the inclusion anywhere,
  then there are chances of the entire list of all tenants being fetched.
- The scopes are too restrictive. There are multiple scopes that can access the resource. In some cases, the scope is
  per account and in some cases, you want a lower ownership.
- Users/Accessors need to have well defined rules for all kinds of memberships.
- Ownership is too limited and cannot be shared or changed on the fly.

In this section, we will try to use the concept of `Group`s to formulate a strategy to fetch a list of objects belonging
to a specific scope.

Let's talk about the different types of scope first. Above, we talk about scopes like `Account` and `Project` which are the
most commonly used across most apps.
There are mini scopes as well which require them to store the scope of the parent.

Proposing a different way of looking at ownership.

What if `Account` or `Project` also leveraged the `Group` model for its memberships instead of maintaining their own memberships?
This would result in an automatic grouping which can be used for authorisation naturally.

Furthering this thought, what if all ownerships or belonging to another object is also done via groups?
Would this result in a more natural and implicit way of deciding authorisation?

```
Any kind of membership or belonging to another object should happen via Groups.
```

### Groups for Accounts and Projects

If the memberships for accounts and projects are stored in Groups, how would it change our access pattern?

When an user tries to access any object belonging to an account, the user first gets added to the account's group.
Once it's added, it should try to unravel the groups it belongs to from the leaf group to the root node.

Similar to the authorisation logic, if it finds any `RoleAccess` rule which allows it to access the object, then
the traversal can stop and the user can have access to it.

We need to define the membership of the object to the account or some other sub account as well.

- Explicit definition of which group the object belongs to.
- Use the user's groups to define the group it belongs to.

This would allow the user to access the object without having any kind of logic on the object models.

For projects, the process is similar.
Let's assume that we want to scope an user only to a specific project.
We add the user to the project's group and the user shouldn't be in the account's group.
From the perspective of groups, the project is also a member of the account's group, i.e they are generic objects
with custom attributes, which we will come to later.

One problem with this approach is to design a system where an user has access to all projects except one. This can be
solved by

- Adding `allow/deny` as an attribute to the `RoleAccess` model.
- Removing the user from the account and instead, adding them to all the individual projects.

Let's not add `deny` to the PRD for now and continue with the 2nd option.

### Groups for child and parent objects

Now let's take an example of a parent object called `Parent` and few children objects named `ChildA`, `ChildB` and `ChildC`.
The children objects exist only in the scope of the parent.

Instead of having a `parent_id` in all the child objects, they can be stored as group members of the `Parent` group
that we can form when the parent is created.

If a parent needs to find the children associated with it, it can perform a query for the `group_type` which will be `Parent`
along with its `id`. It can also have a filter based on the member type if a particular parent owns multiple types of objects.
This will be particularly useful for any of the larger scope groupings.

### Adding an object to a Group

In this section, we will cover how an object can be referred to in a Group, both as a scope and as a member.
Anytime an object needs to be added to a group, it can be done on the fly by directly invoking the Group API or DB write.

Since we are discussing this in the context of `Goiter`, we may need to define more implicit ways of referring to a group.

### Cons of Grouping everything

Few cons of this approach are as follows:

- Ownership is more fluid here. Since there is no direct owner of an object, multiple can own it and that may not be desirable in
  some cases.
- If groups are the sole decider of ownership, it will get bloated quite early on. There are different possible optimisations which
  we will discuss later on.

### Scaling Groups

TBD
