package user

const (
	LIST_OF_USER_FIRST = `
	select 
		user_id,
		full_name,
		user_email,
		birth_date,
		msisdn,
		create_time
	from 
		ws_user 
	where 
		lower(full_name) like '%'||lower($1)||'%' 
	limit 20 
	`
	LIST_OF_USER_SECOND = `
	select 
		user_id,
		full_name,
		user_email,
		birth_date,
		msisdn,
		create_time
	from 
		ws_user	 
	limit 20 
	`
)
