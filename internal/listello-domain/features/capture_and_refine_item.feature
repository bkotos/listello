Feature: Capture and refine an inbox item
  As a user
  I want to capture an inbox item and refine its details before placing it on a list
  So that I can quickly dump work and organize it later

  Scenario: Capturing an inbox item onto a non-inbox list fails
    Given a list named "Next actions" exists
    When the user captures an inbox item "dentist" on the list "Next actions"
    Then the capture should fail with error "can only capture items to the inbox"

  Scenario: Capturing an inbox item
    When the user captures an inbox item "dentist"
    Then a "ItemCaptured" event should have occurred
    And a "ItemCaptured" event should have occurred with the ID of item "dentist"
    And a "ItemCaptured" event should have occurred with the list ID of list "Inbox"
    And the item "dentist" should be outstanding
    And the item "dentist" should have priority "no priority"

  Scenario: Modifying the title of a captured item
    Given a captured item "dentist" exists
    When the owner modifies the title of the item "dentist" to "Schedule dentist"
    Then a "ItemTitleChanged" event should have occurred
    And a "ItemTitleChanged" event should have occurred with the ID of item "Schedule dentist"
    And a "ItemTitleChanged" event should have occurred with title "Schedule dentist"
    And the item "Schedule dentist" should exist

  Scenario: Modifying the description of a captured item
    Given a captured item "Schedule dentist" exists
    When the owner modifies the description of the item "Schedule dentist" to "Call the clinic on Monday"
    Then a "ItemDescriptionChanged" event should have occurred
    And a "ItemDescriptionChanged" event should have occurred with the ID of item "Schedule dentist"
    And a "ItemDescriptionChanged" event should have occurred with description "Call the clinic on Monday"
    And the item "Schedule dentist" should have description "Call the clinic on Monday"

  Scenario: Modifying the due date of a captured item
    Given a captured item "Schedule dentist" exists
    When the owner modifies the due date of the item "Schedule dentist" to "2026-08-03T00:00:00Z"
    Then a "DueDateAddedToItem" event should have occurred
    And a "DueDateAddedToItem" event should have occurred with the ID of item "Schedule dentist"
    And a "DueDateAddedToItem" event should have occurred with due date "2026-08-03T00:00:00Z"
    And the item "Schedule dentist" should be due on "2026-08-03T00:00:00Z"

  Scenario Outline: Accepting a valid ISO due date
    Given a captured item "Schedule dentist" exists
    When the owner modifies the due date of the item "Schedule dentist" to "<due_date>"
    Then a "DueDateAddedToItem" event should have occurred
    And the item "Schedule dentist" should be due on "<due_date>"

    Examples:
      | due_date                 |
      | 2026-08-03T00:00:00Z     |
      | 2026-08-03T15:30:00Z     |
      | 2026-08-03T15:30:00.123Z |
      | 2026-08-03T15:30:00+00:00 |

  Scenario Outline: Rejecting a non-ISO due date
    Given a captured item "Schedule dentist" exists
    When the owner modifies the due date of the item "Schedule dentist" to "<due_date>"
    Then modifying the due date should fail with error "due date must be ISO 8601 format"

    Examples:
      | due_date            |
      | 2026-08-03          |
      | 08/03/2026          |
      | August 3, 2026      |
      | 2026-08-03 00:00:00 |
      | not-a-date          |

  Scenario: Removing the due date of a captured item
    Given a captured item "Schedule dentist" exists
    And the owner modifies the due date of the item "Schedule dentist" to "2026-08-03T00:00:00Z"
    When the owner removes the due date of the item "Schedule dentist"
    Then a "DueDateRemovedFromItem" event should have occurred with the ID of item "Schedule dentist"
    And the item "Schedule dentist" should have no due date

  Scenario: Tagging a captured item
    Given a captured item "Schedule dentist" exists
    When the owner tags the item "Schedule dentist" with "health"
    Then a "TagAddedToItem" event should have occurred
    And a "TagAddedToItem" event should have occurred with the ID of item "Schedule dentist"
    And a "TagAddedToItem" event should have occurred with tag "health"
    And the item "Schedule dentist" should be tagged with "health"

  Scenario Outline: Changing the priority of a captured item
    Given a captured item "Schedule dentist" exists
    When the owner changes the priority of the item "Schedule dentist" to "<priority>"
    Then a "SubtaskPriorityChanged" event should have occurred
    And a "SubtaskPriorityChanged" event should have occurred with the ID of item "Schedule dentist"
    And a "SubtaskPriorityChanged" event should have occurred with priority "<priority>"
    And the item "Schedule dentist" should have priority "<priority>"

    Examples:
      | priority |
      | low      |
      | medium   |
      | high     |

  Scenario: Clearing the priority of a captured item
    Given a captured item "Schedule dentist" exists
    When the owner changes the priority of the item "Schedule dentist" to "high"
    And the owner changes the priority of the item "Schedule dentist" to "no priority"
    Then a "SubtaskPriorityChanged" event should have occurred
    And a "SubtaskPriorityChanged" event should have occurred with priority "no priority"
    And the item "Schedule dentist" should have priority "no priority"

  Scenario: Moving a captured item to a list
    Given a list named "Next actions" exists
    And a captured item "Schedule dentist" exists
    When the owner moves the item "Schedule dentist" to the list "Next actions"
    Then a "ItemMovedToOtherList" event should have occurred
    And a "ItemMovedToOtherList" event should have occurred with the ID of item "Schedule dentist"
    And a "ItemMovedToOtherList" event should have occurred with the list ID of list "Next actions"
    And the item "Schedule dentist" should be on the list "Next actions"

  Scenario: Capturing and refining an item through to a list
    Given a list named "Next actions" exists
    When the user captures an inbox item "dentist"
    And the owner modifies the title of the item "dentist" to "Schedule dentist"
    And the owner modifies the description of the item "Schedule dentist" to "Call the clinic on Monday"
    And the owner modifies the due date of the item "Schedule dentist" to "2026-08-03T00:00:00Z"
    And the owner tags the item "Schedule dentist" with "health"
    And the owner changes the priority of the item "Schedule dentist" to "high"
    And the owner moves the item "Schedule dentist" to the list "Next actions"
    Then a "ItemCaptured" event should have occurred
    And a "ItemTitleChanged" event should have occurred
    And a "ItemDescriptionChanged" event should have occurred
    And a "DueDateAddedToItem" event should have occurred
    And a "TagAddedToItem" event should have occurred
    And a "SubtaskPriorityChanged" event should have occurred
    And a "ItemMovedToOtherList" event should have occurred
    And the item "Schedule dentist" should have description "Call the clinic on Monday"
    And the item "Schedule dentist" should be due on "2026-08-03T00:00:00Z"
    And the item "Schedule dentist" should be tagged with "health"
    And the item "Schedule dentist" should have priority "high"
    And the item "Schedule dentist" should be on the list "Next actions"
    And the item "Schedule dentist" should be outstanding
