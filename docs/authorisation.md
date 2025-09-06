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

### Note

- Have separate rules to govern for owners of objects
- Should we have separate group models for different elements or a single group model for all the elements
  with an identifier?
