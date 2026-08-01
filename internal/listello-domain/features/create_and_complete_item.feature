@wip
Feature: Create and complete a list item
  As a user
  I want to create a list, define an item on it, and complete that item
  So that I can track work from empty list through to done

  Scenario: Creating a list
    When the user creates a list named "Next actions"
    Then a "List created" event should have occurred
    And the list "Next actions" should exist

  Scenario: Defining an item on a list
    Given a list named "Next actions" exists
    When the user defines an item titled "Buy milk" on the list "Next actions"
    Then a "Item defined" event should have occurred
    And the item "Buy milk" should be on the list "Next actions"
    And the item "Buy milk" should be incomplete

  Scenario: Completing an item
    Given a list named "Next actions" exists
    And an incomplete item titled "Buy milk" exists on the list "Next actions"
    When the owner completes the item "Buy milk"
    Then a "Item completed" event should have occurred
    And the item "Buy milk" should be complete

  Scenario: Completing work from a new list through to done
    When the user creates a list named "Next actions"
    And the user defines an item titled "Buy milk" on the list "Next actions"
    And the owner completes the item "Buy milk"
    Then a "List created" event should have occurred
    And a "Item defined" event should have occurred
    And a "Item completed" event should have occurred
    And the list "Next actions" should exist
    And the item "Buy milk" should be on the list "Next actions"
    And the item "Buy milk" should be complete
