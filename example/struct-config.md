@startuml

general : Global configuration of your dashboard
projects : List of your projects
services: Configurations of every services you want to use
widgets: List of widgets you want to display
row: Create a row which contains columns
col: Create a column which contains widgets
size: Size of the column (T-shirt sizes or number 0-12)
elements: Your actual widgets and their configuration for the current column


general-->projects
projects-->services
projects--->widgets
widgets-->row
row-->col
col-->size
col--->elements

@enduml

