Feature: Delete a comment on an item
  As a user
  I want to delete a comment I authored
  So that I can remove a note I no longer need

  Only the comment author may delete a comment. Delete records who deleted it and when.

  Scenario: The author deletes a comment on an item
    Given a list named "Next actions" exists
    And an outstanding item titled "Buy milk" exists on the list "Next actions"
    And a user named "Alex" exists
    And the user "Alex" comments "Need 2%" on the item "Buy milk"
    When the user "Alex" deletes the comment "Need 2%" on the item "Buy milk"
    Then a "ItemCommentDeleted" event should have occurred
    And a "ItemCommentDeleted" event should have occurred with the ID of item "Buy milk"
    And a "ItemCommentDeleted" event should have occurred with the ID of comment "Need 2%" on the item "Buy milk"
    And a "ItemCommentDeleted" event should have occurred with the deleter "Alex"
    And a "ItemCommentDeleted" event should have occurred with a deleted-at ISO date time
    And the item "Buy milk" should not have a comment "Need 2%"

  Scenario: A non-author cannot delete a comment
    Given a list named "Next actions" exists
    And an outstanding item titled "Buy milk" exists on the list "Next actions"
    And a user named "Alex" exists
    And a user named "Sam" exists
    And the user "Alex" comments "Need 2%" on the item "Buy milk"
    When the user "Sam" deletes the comment "Need 2%" on the item "Buy milk"
    Then deleting the comment should fail with error "only the comment author can delete the comment"
    And the item "Buy milk" should have a comment "Need 2%"
    And a "ItemCommentDeleted" event should not have occurred

  Scenario: Deleting an unknown comment is rejected
    Given a list named "Next actions" exists
    And an outstanding item titled "Buy milk" exists on the list "Next actions"
    And a user named "Alex" exists
    When the user "Alex" deletes a comment with ID "CM_unknown" on the item "Buy milk"
    Then deleting the comment should fail with error "comment not found"
    And a "ItemCommentDeleted" event should not have occurred
