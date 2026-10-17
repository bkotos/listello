Feature: Delete a comment on an item
  As a user
  I want to delete a comment on an item
  So that I can remove a note I no longer need

  Scenario: Deleting a comment on an item
    Given a list named "Next actions" exists
    And an outstanding item titled "Buy milk" exists on the list "Next actions"
    And a user named "Alex" exists
    And the user "Alex" comments "Need 2%" on the item "Buy milk"
    When the user deletes the comment "Need 2%" on the item "Buy milk"
    Then a "ItemCommentDeleted" event should have occurred
    And a "ItemCommentDeleted" event should have occurred with the ID of item "Buy milk"
    And a "ItemCommentDeleted" event should have occurred with the ID of comment "Need 2%" on the item "Buy milk"
    And the item "Buy milk" should not have a comment "Need 2%"
