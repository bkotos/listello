@wip
Feature: Set up personal productivity
  As a user
  I want to create a space and myself as a user, then have the system create an inbox and assign the space, then create my first list
  So that I can start capturing and organizing work

  Scenario: Creating a space
    When the user creates a space named "Personal"
    Then a "SpaceCreated" event should have occurred
    And the space "Personal" should exist

  Scenario: Creating themselves as a user
    When the user creates themselves as a user named "Alex"
    Then a "UserCreated" event should have occurred
    And the user "Alex" should exist

  Scenario: Creating an inbox
    Given a space named "Personal" exists
    When the system creates an inbox
    Then a "InboxCreated" event should have occurred
    And the list "Inbox" should exist
    And the inbox should be attached to the space "Personal"

  Scenario: Assigning a space to a user
    Given a space named "Personal" exists
    And a user named "Alex" exists
    When the system assigns the space "Personal" to the user "Alex"
    Then a "SpaceAssignedToUser" event should have occurred
    And the space "Personal" should be assigned to the user "Alex"

  Scenario: Creating their first list
    Given a space named "Personal" exists
    And a user named "Alex" exists
    And the system assigns the space "Personal" to the user "Alex"
    And the system created an inbox
    And the user has no lists other than the inbox
    When the user creates their first list named "Next actions"
    Then a "ListCreated" event should have occurred
    And a "FirstListCreated" event should have occurred
    And the list "Next actions" should exist

  Scenario: Setting up personal productivity through to a first list
    When the user creates a space named "Personal"
    And the user creates themselves as a user named "Alex"
    And the system creates an inbox
    And the system assigns the space "Personal" to the user "Alex"
    And the user creates their first list named "Next actions"
    Then a "SpaceCreated" event should have occurred
    And a "UserCreated" event should have occurred
    And a "InboxCreated" event should have occurred
    And a "SpaceAssignedToUser" event should have occurred
    And a "ListCreated" event should have occurred
    And a "FirstListCreated" event should have occurred
    And the space "Personal" should exist
    And the user "Alex" should exist
    And the list "Inbox" should exist
    And the space "Personal" should be assigned to the user "Alex"
    And the list "Next actions" should exist
