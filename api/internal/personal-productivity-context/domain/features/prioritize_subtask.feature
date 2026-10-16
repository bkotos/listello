Feature: Prioritize a subtask on an item
  As an owner
  I want to change the priority of a subtask on an item
  So that I can rank a smaller step without changing the item's priority

  Scenario: Changing the priority of a subtask on an item
    Given a list named "Next actions" exists
    And an outstanding item titled "Buy milk" exists on the list "Next actions"
    And the owner adds a subtask "Get 2%" to the item "Buy milk"
    When the owner changes the priority of the subtask "Get 2%" on the item "Buy milk" to "high"
    Then a "SubtaskPriorityChanged" event should have occurred
    And a "SubtaskPriorityChanged" event should have occurred with the ID of item "Buy milk"
    And a "SubtaskPriorityChanged" event should have occurred with the ID of subtask "Get 2%" on the item "Buy milk"
    And a "SubtaskPriorityChanged" event should have occurred with priority "high"
    And the subtask "Get 2%" on the item "Buy milk" should have priority "high"
