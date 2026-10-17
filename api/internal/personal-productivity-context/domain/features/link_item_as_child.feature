Feature: Link an item as a child of another item
  As an owner
  I want to link an existing item as a child of another item
  So that I can nest work without converting the child into a subtask

  Scenario: Linking an item as a child of another item
    Given a list named "Next actions" exists
    And an outstanding item titled "Plan trip" exists on the list "Next actions"
    And an outstanding item titled "Book flights" exists on the list "Next actions"
    When the owner links the item "Book flights" as a child of the item "Plan trip"
    Then a "ItemLinkedAsChildOfItem" event should have occurred
    And a "ItemLinkedAsChildOfItem" event should have occurred with the ID of item "Book flights"
    And a "ItemLinkedAsChildOfItem" event should have occurred with the parent ID of item "Plan trip"
    And the item "Book flights" should be a child of the item "Plan trip"
    And the item "Book flights" should have an ID prefixed with "IT_"
    And the item "Plan trip" should not have a subtask "Book flights"

  Scenario: Linking multiple items as children of the same parent
    Given a list named "Next actions" exists
    And an outstanding item titled "Plan trip" exists on the list "Next actions"
    And an outstanding item titled "Book flights" exists on the list "Next actions"
    And an outstanding item titled "Reserve hotel" exists on the list "Next actions"
    When the owner links the item "Book flights" as a child of the item "Plan trip"
    And the owner links the item "Reserve hotel" as a child of the item "Plan trip"
    Then the item "Book flights" should be a child of the item "Plan trip"
    And the item "Reserve hotel" should be a child of the item "Plan trip"
    And the item "Book flights" should have an ID prefixed with "IT_"
    And the item "Reserve hotel" should have an ID prefixed with "IT_"
    And the item "Plan trip" should not have a subtask "Book flights"
    And the item "Plan trip" should not have a subtask "Reserve hotel"

  Scenario: Linking an item as a child of an item on another list
    Given a list named "Work" exists
    And a list named "Personal" exists
    And an outstanding item titled "Plan trip" exists on the list "Personal"
    And an outstanding item titled "Book flights" exists on the list "Work"
    When the owner links the item "Book flights" as a child of the item "Plan trip"
    Then the item "Book flights" should be a child of the item "Plan trip"
    And the item "Book flights" should have an ID prefixed with "IT_"
    And the item "Plan trip" should not have a subtask "Book flights"

  Scenario: Linking an item as a grandchild fails
    Given a list named "Next actions" exists
    And an outstanding item titled "Plan trip" exists on the list "Next actions"
    And an outstanding item titled "Book flights" exists on the list "Next actions"
    And an outstanding item titled "Compare airlines" exists on the list "Next actions"
    And the owner links the item "Book flights" as a child of the item "Plan trip"
    When the owner links the item "Compare airlines" as a child of the item "Book flights"
    Then linking as a child should fail with error "cannot link as child of an item that already has a parent"
