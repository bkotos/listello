Feature: Delete a subtask on an item
  As an owner
  I want to delete a subtask on an item
  So that I can remove a smaller step I no longer need

  Scenario: Deleting a subtask on an item
    Given a list named "Next actions" exists
    And an outstanding item titled "Buy milk" exists on the list "Next actions"
    And the owner adds a subtask "Get 2%" to the item "Buy milk"
    When the owner deletes the subtask "Get 2%" on the item "Buy milk"
    Then a "SubtaskDeletedOnItem" event should have occurred
    And a "SubtaskDeletedOnItem" event should have occurred with the ID of item "Buy milk"
    And a "SubtaskDeletedOnItem" event should have occurred with the ID of subtask "Get 2%" on the item "Buy milk"
    And the item "Buy milk" should not have a subtask "Get 2%"
