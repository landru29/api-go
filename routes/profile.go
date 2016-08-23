package routes

// Profile define a user profile
type Profile struct {
    ID        string `json:"id"`
    Email     string `json:"email"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Hd        string `json:"hd"`
    Locale    string `json:"locale"`
    Name      string `json:"name"`
    Picture   struct {
        Data struct {
            URL string `json:"url"`
        } `json:"data"`
    } `json:"picture"`
}
