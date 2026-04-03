package server
import ("encoding/json";"log";"net/http";"github.com/stockyard-dev/stockyard-pasture/internal/store")
type Server struct{db *store.DB;mux *http.ServeMux;limits Limits}
func New(db *store.DB,limits Limits)*Server{s:=&Server{db:db,mux:http.NewServeMux(),limits:limits}
s.mux.HandleFunc("GET /api/posts",s.list);s.mux.HandleFunc("POST /api/posts",s.create);s.mux.HandleFunc("GET /api/posts/{id}",s.get);s.mux.HandleFunc("PUT /api/posts/{id}",s.update);s.mux.HandleFunc("DELETE /api/posts/{id}",s.del)
s.mux.HandleFunc("POST /api/posts/{id}/publish",s.publish)
s.mux.HandleFunc("GET /api/stats",s.stats);s.mux.HandleFunc("GET /api/health",s.health)
s.mux.HandleFunc("GET /ui",s.dashboard);s.mux.HandleFunc("GET /ui/",s.dashboard);s.mux.HandleFunc("GET /",s.root);return s}
func(s *Server)ServeHTTP(w http.ResponseWriter,r *http.Request){s.mux.ServeHTTP(w,r)}
func wj(w http.ResponseWriter,c int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(c);json.NewEncoder(w).Encode(v)}
func we(w http.ResponseWriter,c int,m string){wj(w,c,map[string]string{"error":m})}
func(s *Server)root(w http.ResponseWriter,r *http.Request){if r.URL.Path!="/"{http.NotFound(w,r);return};http.Redirect(w,r,"/ui",302)}
func(s *Server)list(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"posts":oe(s.db.List(r.URL.Query().Get("status")))})}
func(s *Server)create(w http.ResponseWriter,r *http.Request){var p store.Post;json.NewDecoder(r.Body).Decode(&p);if p.Title==""{we(w,400,"title required");return};s.db.Create(&p);wj(w,201,s.db.Get(p.ID))}
func(s *Server)get(w http.ResponseWriter,r *http.Request){p:=s.db.Get(r.PathValue("id"));if p==nil{we(w,404,"not found");return};wj(w,200,p)}
func(s *Server)update(w http.ResponseWriter,r *http.Request){id:=r.PathValue("id");ex:=s.db.Get(id);if ex==nil{we(w,404,"not found");return};var p store.Post;json.NewDecoder(r.Body).Decode(&p);if p.Title==""{p.Title=ex.Title};if p.Status==""{p.Status=ex.Status};s.db.Update(id,&p);wj(w,200,s.db.Get(id))}
func(s *Server)del(w http.ResponseWriter,r *http.Request){s.db.Delete(r.PathValue("id"));wj(w,200,map[string]string{"deleted":"ok"})}
func(s *Server)publish(w http.ResponseWriter,r *http.Request){s.db.Publish(r.PathValue("id"));wj(w,200,s.db.Get(r.PathValue("id")))}
func(s *Server)stats(w http.ResponseWriter,r *http.Request){wj(w,200,s.db.Stats())}
func(s *Server)health(w http.ResponseWriter,r *http.Request){st:=s.db.Stats();wj(w,200,map[string]any{"status":"ok","service":"pasture","posts":st.Total,"scheduled":st.Scheduled})}
func oe[T any](s []T)[]T{if s==nil{return[]T{}};return s}
func init(){log.SetFlags(log.LstdFlags|log.Lshortfile)}
