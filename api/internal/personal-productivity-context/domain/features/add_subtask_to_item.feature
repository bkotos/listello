Feature: Add a subtask to an item
  As an owner
  I want to add a subtask to an item
  So that I can break work into smaller steps

  Scenario: Adding a subtask to an item
    Given a list named "Next actions" exists
    And an outstanding item titled "Buy milk" exists on the list "Next actions"
    When the owner adds a subtask "Get 2%" to the item "Buy milk"
    Then a "SubtaskAddedToItem" event should have occurred
    And a "SubtaskAddedToItem" event should have occurred with the ID of item "Buy milk"
    And a "SubtaskAddedToItem" event should have occurred with the ID of subtask "Get 2%" on the item "Buy milk"
    And the item "Buy milk" should have a subtask "Get 2%"
    And the subtask "Get 2%" on the item "Buy milk" should have an ID prefixed with "ST_"
    And the subtask "Get 2%" on the item "Buy milk" should not be completed
