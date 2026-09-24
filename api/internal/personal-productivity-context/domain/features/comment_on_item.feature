Feature: Comment on an item
  As a user
  I want to comment on an item
  So that I can record notes against the work

  Scenario: Commenting on an item
    Given a list named "Next actions" exists
    And an outstanding item titled "Buy milk" exists on the list "Next actions"
    When the user comments "Need 2%" on the item "Buy milk"
    Then a "ItemCommentedOn" event should have occurred
    And a "ItemCommentedOn" event should have occurred with the ID of item "Buy milk"
    And a "ItemCommentedOn" event should have occurred with body "Need 2%"
    And the item "Buy milk" should have a comment "Need 2%"
    And the comment "Need 2%" on the item "Buy milk" should have an ID prefixed with "CM_"
