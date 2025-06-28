# LOG-BOOK PROJECT
## OPEN ISSUES LIST
---------------------------------------------------------------------------------------------------
### improved_user_input
- [ ] 2025-06-27: See if there are any other uses for inputState.statusMsg

- [ ] 2025-06-27: Add Ctrl-D safety to the project selection and display screens

- [ ] 2025-06-27: If project selection is incorrect, immediately bring user to the helpful list of
                  available projects

- [ ] 2025-06-27: Add a more helpful error message when user selected an out-of-range project, or
                  nothing at all. Rebuild the project list screen, and add a STATUS MESSAGE
                  section, similar to the main menu solution
         
### validation_before_write
- [ ] 2025-06-21: Create validation function before writing to SQL database

### reporting_feature
- [ ] 2025-06-24: Add export to Excel (or csv if I get lazy)
- [ ] 2025-06-24: Add summary reports
                    -- project summaries
                    -- weekly summaries
                    -- monthly summaries

## Completed Items
- [x] 2025-06-21: Cancel input should zero out logData

- [x] 2025-06-21: Remove numbering from main menu - use letter hints only

- [x] 2025-06-21: display full log entry after each input

- [x] 2025-06-21: Remove the "Input:" prompt after a selection has been made. Replace with Ctrl-D
                  hint for exiting

- [x] 2025-06-21: Also remove the "Input:" prompt and replace with Ctrl-D hint for exiting when in
                  EDIT MODE
- [x] 2025-06-21: Change 'Q' for quit to Ctrl-D to cancel data entry (since "Q" cannot be used in
                  text fields)

- [x] 2025-06-22: Change date display to show immediately after entry (with 00:00 for time)

- [x] 2025-06-22: Add logic to duration calc so that if beginning or ending times are zero, don't
                  calculate - tricky, since 00:00 == midnight

- [x] 2025-06-21: In Edit mode, if no date has been detected before entering start or end time
                  bring user to the date function. 

- [x] 2025-06-21: Display "INPUT MODE" for initial input to make it consistent with "EDIT MODE"

- [x] 2025-06-21: Change Ctrl-D behavior when at menus (main menu, edit menu, input menu) so that
                  it is ignored (i.e., just continue the loop)

- [x] 2025-06-27: `fmt.Printf("Successfully saved log entry with ID: %d\n", insertId)` line does
                  not display because it is immediately overwritten by the new display. Let's add
                  a string field to InputState statusMsg to hold values that we want to display
                  after the menu display 

