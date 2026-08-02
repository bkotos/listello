Feature: Create and complete a list item
  As a user
  I want to create a list, define an item on it, and complete that item
  So that I can track work from empty list through to done

  Scenario: Creating a list named Inbox fails
    When the user creates a list named "Inbox"
    Then creating the list should fail with error "cannot create a list named Inbox"

  Scenario: Creating a list
    When the user creates a list named "Next actions"
    Then a "ListCreated" event should have occurred
    And the list "Next actions" should exist

  Scenario: Defining an item on the inbox fails
    Given an inbox list exists
    When the user defines an item titled "Buy milk" on the list "Inbox"
    Then defining the item should fail with error "can only capture items on inbox lists, not define them"

  Scenario: Defining an item on a list
    Given a list named "Next actions" exists
    When the user defines an item titled "Buy milk" on the list "Next actions"
    Then a "ItemDefined" event should have occurred
    And the item "Buy milk" should be on the list "Next actions"
    And the item "Buy milk" should be outstanding

  Scenario: Completing an item
    Given a list named "Next actions" exists
    And an outstanding item titled "Buy milk" exists on the list "Next actions"
    When the owner completes the item "Buy milk"
    Then a "ItemCompleted" event should have occurred
    And the item "Buy milk" should be complete

  Scenario: Completing work from a new list through to done
    When the user creates a list named "Next actions"
    And the user defines an item titled "Buy milk" on the list "Next actions"
    And the owner completes the item "Buy milk"
    Then a "ListCreated" event should have occurred
    And a "ItemDefined" event should have occurred
    And a "ItemCompleted" event should have occurred
    And the list "Next actions" should exist
    And the item "Buy milk" should be on the list "Next actions"
    And the item "Buy milk" should be complete
