# Adding Your Own Models 

Authorization based on Resource Type And Belonging is SOLELY the responsibility of the models
Roles in the default implementations (Postgres, DynamoDB) only allow owners/creators of a space to create shares,
and invites to shares.

This can be adjusted with diferent logic in the model as each interface method gets authentication context


Routes to Promote Roles Have not been added but plan to be in the future