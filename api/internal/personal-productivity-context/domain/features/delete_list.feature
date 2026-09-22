Feature: Delete a list
  As a user
  I want to delete a list
  So that I can remove lists I no longer need

  Scenario: Deleting a list
    Given a list named "Next actions" exists
    When the user deletes the list "Next actions"
    Then a "ListDeleted" event should have occurred with the ID of list "Next actions"
    And a "ListDeleted" event should have occurred with the list "Next actions"
    And the list "Next actions" should not exist
