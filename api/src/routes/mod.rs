use actix_web::web;
use crate::handlers;

pub fn configure_routes(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/api/todos")
            .route("", web::get().to(handlers::list_todos))
            .route("", web::post().to(handlers::create_todo))
            .route("/{id}", web::get().to(handlers::get_todo))
            .route("/{id}", web::put().to(handlers::update_todo))
            .route("/{id}", web::delete().to(handlers::delete_todo))
    );
}
