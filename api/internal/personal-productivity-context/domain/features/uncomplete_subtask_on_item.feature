Feature: Uncomplete a subtask on an item
  As an owner
  I want to uncomplete a subtask on an item
  So that I can mark a smaller step as not done

  Scenario: Uncompleting a subtask on an item
    Given a list named "Next actions" exists
    And an outstanding item titled "Buy milk" exists on the list "Next actions"
    And the owner adds a subtask "Get 2%" to the item "Buy milk"
    And the owner completes the subtask "Get 2%" on the item "Buy milk"
    When the owner uncompletes the subtask "Get 2%" on the item "Buy milk"
    Then a "SubtaskUncompleted" event should have occurred
    And a "SubtaskUncompleted" event should have occurred with the ID of item "Buy milk"
    And a "SubtaskUncompleted" event should have occurred with the ID of subtask "Get 2%" on the item "Buy milk"
    And the subtask "Get 2%" on the item "Buy milk" should not be completed
