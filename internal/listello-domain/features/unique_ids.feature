Feature: Unique IDs for lists and items
  As a user
  I want lists and items to receive unique IDs when created
  So that each can be referenced independently of its name or title

  Scenario: Creating a list assigns an ID prefixed with LS_
    When the user creates a list named "Next actions"
    Then the list "Next actions" should have an ID prefixed with "LS_"

  Scenario: Creating an item assigns an ID prefixed with IT_
    Given a list named "Next actions" exists
    When the user defines an item titled "Buy milk" on the list "Next actions"
    Then the item "Buy milk" should have an ID prefixed with "IT_"

  Scenario: Creating lists assigns distinct IDs
    When the user creates a list named "Next actions"
    And the user creates a list named "Someday"
    Then the lists "Next actions" and "Someday" should have different IDs

  Scenario: Creating items assigns distinct IDs
    Given a list named "Next actions" exists
    When the user defines an item titled "Buy milk" on the list "Next actions"
    And the user defines an item titled "Buy eggs" on the list "Next actions"
    Then the items "Buy milk" and "Buy eggs" should have different IDs
