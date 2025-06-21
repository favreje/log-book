# LOG-BOOK PROJECT
## OPEN ISSUES LIST
---------------------------------------------------------------------------------------------------
- [ ] 2025-06-21: Also remove the "Input:" prompt and replace with Ctrl-D hint for exiting when in
                  EDIT MODE
  
- [ ] 2025-06-21: Change 'Q' for quit to Ctrl-D to cancel data entry (since "Q" cannot be used in
                  text fields)
  
- [ ] 2025-06-21: Create validation function before writing to SQL database
  
- [ ] 2025-06-21: Display "INPUT MODE" for initial input to make it consistent with "EDIT MODE"
  
- [ ] 2025-06-21: Change Ctrl-D behavior when at menus (main menu, edit menu, input menu) so that it
                  is ignored (i.e., just continue the loop)
  
- [ ] 2025-06-21: In Edit mode, if no date has been detected before entering start or end time bring
                  user to the date function. This is to avoid having a "01/01/01" zero value display
                  for the date portion of start and end times.
  
- [ ] 2025-06-21: Change boundaryType to an enum rather than using strings (to keep it classy)


## Completed Items
- [x] 2025-06-21: Cancel input should zero out logData
- [x] 2025-06-21: Remove numbering from main menu - use letter hints only
- [x] 2025-06-21: display full log entry after each input
- [x] 2025-06-21: Remove the "Input:" prompt after a selection has been made. Replace with Ctrl-D hint for exiting
