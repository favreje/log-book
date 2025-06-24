# LOG-BOOK PROJECT
## OPEN ISSUES LIST
---------------------------------------------------------------------------------------------------
### improved_user_input
- [ ] 2025-06-21: Display "INPUT MODE" for initial input to make it consistent with "EDIT MODE"

- [ ] 2025-06-21: Change boundaryType to an enum rather than using strings (to keep it classy)
  
- [ ] 2025-06-21: Change Ctrl-D behavior when at menus (main menu, edit menu, input menu) so that it
                  is ignored (i.e., just continue the loop)
  

### validation_before_write
- [ ] 2025-06-21: Create validation function before writing to SQL database


## Completed Items
- [x] 2025-06-21: Cancel input should zero out logData
- [x] 2025-06-21: Remove numbering from main menu - use letter hints only
- [x] 2025-06-21: display full log entry after each input
- [x] 2025-06-21: Remove the "Input:" prompt after a selection has been made. Replace with Ctrl-D hint for exiting
- [x] 2025-06-21: Also remove the "Input:" prompt and replace with Ctrl-D hint for exiting when in EDIT MODE
- [x] 2025-06-21: Change 'Q' for quit to Ctrl-D to cancel data entry (since "Q" cannot be used in text fields)
- [x] 2025-06-22: Change date display to show immediately after entry (with 00:00 for time)
- [x] 2025-06-22: Add logic to duration calc so that if beginning or ending times are zero, don't calculate - tricky, since 00:00 == midnight
- [x] 2025-06-21: In Edit mode, if no date has been detected before entering start or end time bring user to the date function. 
