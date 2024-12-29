package main

import (
    "encoding/json"
	    "fmt"
		    "log"
			    "net/http"
				    "os/exec"
					    "strings"

						    "github.com/gorilla/mux"
							)

							// Repository represents a Docker repository
							type Repository struct {
							    Name        string   `json:"name"`
								    Description string   `json:"description"`
									    StarCount   int      `json:"star_count"`
										    PullCount   int      `json:"pull_count"`
											    LastUpdated string   `json:"last_updated"`
												    Tags        []string `json:"tags"`
													}

													// Settings represents the settings for accelerators
													type Settings struct {
													    Accelerators []string `json:"accelerators"`
														}

														var settings Settings

														func main() {
														    router := mux.NewRouter()

															    router.HandleFunc("/api/public_repositories", GetPublicRepositories).Methods("GET")
																    router.HandleFunc("/api/search", SearchImages).Methods("GET")
																	    router.HandleFunc("/api/pull", PullImage).Methods("POST")
																		    router.HandleFunc("/api/repository/{name}", GetRepositoryDetails).Methods("GET")
																			    router.HandleFunc("/api/settings", UpdateSettings).Methods("POST")
																				    router.HandleFunc("/api/settings", GetSettings).Methods("GET")

																					    // Serve frontend files
																						    router.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend/")))

																							    log.Println("Starting server on :8080")
																								    log.Fatal(http.ListenAndServe(":8080", router))
																									}

																									// GetPublicRepositories fetches the list of public repositories from Docker Hub
																									func GetPublicRepositories(w http.ResponseWriter, r *http.Request) {
																									    log.Println("Fetching public repositories")

																										    url := "https://hub.docker.com/v2/repositories/library/?page_size=10"
																											    if len(settings.Accelerators) > 0 {
																												        url = strings.Replace(url, "https://hub.docker.com", settings.Accelerators[0], 1)
																														    }

																															    response, err := http.Get(url)
																																    if err != nil {
																																	        log.Printf("Error fetching public repositories: %v\n", err)
																																			        http.Error(w, err.Error(), http.StatusInternalServerError)
																																					        return
																																							    }
																																								    defer response.Body.Close()

																																									    var data map[string]interface{}
																																										    if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
																																											        log.Printf("Error decoding response: %v\n", err)
																																													        http.Error(w, err.Error(), http.StatusInternalServerError)
																																															        return
																																																	    }

																																																		    repositories, ok := data["results"].([]interface{})
																																																			    if !ok {
																																																				        log.Printf("Error: invalid response format\n")
																																																						        http.Error(w, "Invalid response format", http.StatusInternalServerError)
																																																								        return
																																																										    }

																																																											    var result []Repository
																																																												    for _, repo := range repositories {
																																																													        repoData, ok := repo.(map[string]interface{})
																																																															        if !ok {
																																																																	            log.Printf("Error: invalid repository data format\n")
																																																																				            continue
																																																																							        }
																																																																									        name, nameOk := repoData["repo_name"].(string)
																																																																											        description, descriptionOk := repoData["short_description"].(string)
																																																																													        starCount, starCountOk := repoData["star_count"].(float64)
																																																																															        pullCount, pullCountOk := repoData["pull_count"].(float64)

																																																																																	        // 确保所有字段都存在
																																																																																			        if !nameOk || !descriptionOk || !starCountOk || !pullCountOk {
																																																																																					            log.Printf("Error: missing fields in repository data\n")
																																																																																								            continue
																																																																																											        }

																																																																																													        repository := Repository{
																																																																																															            Name:        name,
																																																																																																		            Description: description,
																																																																																																					            StarCount:   int(starCount),
																																																																																																								            PullCount:   int(pullCount),
																																																																																																											            LastUpdated: "", // 原始数据中没有 last_updated 字段
																																																																																																														            Tags:        []string{},
																																																																																																																	        }
																																																																																																																			        result = append(result, repository)
																																																																																																																					    }

																																																																																																																						    w.Header().Set("Content-Type", "application/json")
																																																																																																																							    json.NewEncoder(w).Encode(result)
																																																																																																																								    log.Println("Fetched public repositories successfully")
																																																																																																																									}

																																																																																																																									// SearchImages searches for Docker images on Docker Hub
																																																																																																																									func SearchImages(w http.ResponseWriter, r *http.Request) {
																																																																																																																									    query := r.URL.Query().Get("query")
																																																																																																																										    log.Printf("Searching images with query: %s\n", query)

																																																																																																																											    url := fmt.Sprintf("https://hub.docker.com/v2/search/repositories/?query=%s&page_size=10", query)
																																																																																																																												    if len(settings.Accelerators) > 0 {
																																																																																																																													        url = strings.Replace(url, "https://hub.docker.com", settings.Accelerators[0], 1)
																																																																																																																															    }

																																																																																																																																    response, err := http.Get(url)
																																																																																																																																	    if err != nil {
																																																																																																																																		        log.Printf("Error searching images: %v\n", err)
																																																																																																																																				        http.Error(w, err.Error(), http.StatusInternalServerError)
																																																																																																																																						        return
																																																																																																																																								    }
																																																																																																																																									    defer response.Body.Close()

																																																																																																																																										    var data map[string]interface{}
																																																																																																																																											    if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
																																																																																																																																												        log.Printf("Error decoding response: %v\n", err)
																																																																																																																																														        http.Error(w, err.Error(), http.StatusInternalServerError)
																																																																																																																																																        return
																																																																																																																																																		    }

																																																																																																																																																			    repositories, ok := data["results"].([]interface{})
																																																																																																																																																				    if !ok {
																																																																																																																																																					        log.Printf("Error: invalid response format\n")
																																																																																																																																																							        http.Error(w, "Invalid response format", http.StatusInternalServerError)
																																																																																																																																																									        return
																																																																																																																																																											    }

																																																																																																																																																												    var result []Repository
																																																																																																																																																													    for _, repo := range repositories {
																																																																																																																																																														        repoData, ok := repo.(map[string]interface{})
																																																																																																																																																																        if !ok {
																																																																																																																																																																		            log.Printf("Error: invalid repository data format\n")
																																																																																																																																																																					            continue
																																																																																																																																																																								        }
																																																																																																																																																																										        name, nameOk := repoData["repo_name"].(string)
																																																																																																																																																																												        description, descriptionOk := repoData["short_description"].(string)
																																																																																																																																																																														        starCount, starCountOk := repoData["star_count"].(float64)
																																																																																																																																																																																        pullCount, pullCountOk := repoData["pull_count"].(float64)

																																																																																																																																																																																		        // 确保所有字段都存在
																																																																																																																																																																																				        if !nameOk || !descriptionOk || !starCountOk || !pullCountOk {
																																																																																																																																																																																						            log.Printf("Error: missing fields in repository data\n")
																																																																																																																																																																																									            continue
																																																																																																																																																																																												        }

																																																																																																																																																																																														        repository := Repository{
																																																																																																																																																																																																            Name:        name,
																																																																																																																																																																																																			            Description: description,
																																																																																																																																																																																																						            StarCount:   int(starCount),
																																																																																																																																																																																																									            PullCount:   int(pullCount),
																																																																																																																																																																																																												            LastUpdated: "", // 原始数据中没有 last_updated 字段
																																																																																																																																																																																																															            Tags:        []string{},
																																																																																																																																																																																																																		        }
																																																																																																																																																																																																																				        result = append(result, repository)
																																																																																																																																																																																																																						    }

																																																																																																																																																																																																																							    w.Header().Set("Content-Type", "application/json")
																																																																																																																																																																																																																								    json.NewEncoder(w).Encode(result)
																																																																																																																																																																																																																									    log.Println("Search images completed successfully")
																																																																																																																																																																																																																										}

																																																																																																																																																																																																																										// GetRepositoryDetails fetches the details of a repository, including tags
																																																																																																																																																																																																																										func GetRepositoryDetails(w http.ResponseWriter, r *http.Request) {
																																																																																																																																																																																																																										    vars := mux.Vars(r)
																																																																																																																																																																																																																											    name := vars["name"]
																																																																																																																																																																																																																												    log.Printf("Fetching details for repository: %s\n", name)

																																																																																																																																																																																																																													    url := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/tags/", name)
																																																																																																																																																																																																																														    if len(settings.Accelerators) > 0 {
																																																																																																																																																																																																																															        url = strings.Replace(url, "https://hub.docker.com", settings.Accelerators[0], 1)
																																																																																																																																																																																																																																	    }

																																																																																																																																																																																																																																		    response, err := http.Get(url)
																																																																																																																																																																																																																																			    if err != nil {
																																																																																																																																																																																																																																				        log.Printf("Error fetching repository details: %v\n", err)
																																																																																																																																																																																																																																						        http.Error(w, err.Error(), http.StatusInternalServerError)
																																																																																																																																																																																																																																								        return
																																																																																																																																																																																																																																										    }
																																																																																																																																																																																																																																											    defer response.Body.Close()

																																																																																																																																																																																																																																												    var data map[string]interface{}
																																																																																																																																																																																																																																													    if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
																																																																																																																																																																																																																																														        log.Printf("Error decoding response: %v\n", err)
																																																																																																																																																																																																																																																        http.Error(w, err.Error(), http.StatusInternalServerError)
																																																																																																																																																																																																																																																		        return
																																																																																																																																																																																																																																																				    }

																																																																																																																																																																																																																																																					    tagsData, ok := data["results"].([]interface{})
																																																																																																																																																																																																																																																						    if !ok {
																																																																																																																																																																																																																																																							        log.Printf("Error: invalid response format\n")
																																																																																																																																																																																																																																																									        http.Error(w, "Invalid response format", http.StatusInternalServerError)
																																																																																																																																																																																																																																																											        return
																																																																																																																																																																																																																																																													    }

																																																																																																																																																																																																																																																														    var tags []string
																																																																																																																																																																																																																																																															    for _, tag := range tagsData {
																																																																																																																																																																																																																																																																        tagData, ok := tag.(map[string]interface{})
																																																																																																																																																																																																																																																																		        if !ok {
																																																																																																																																																																																																																																																																				            log.Printf("Error: invalid tag data format\n")
																																																																																																																																																																																																																																																																							            continue
																																																																																																																																																																																																																																																																										        }
																																																																																																																																																																																																																																																																												        tagName, tagNameOk := tagData["name"].(string)
																																																																																																																																																																																																																																																																														        if !tagNameOk {
																																																																																																																																																																																																																																																																																            log.Printf("Error: missing tag name in tag data\n")
																																																																																																																																																																																																																																																																																			            continue
																																																																																																																																																																																																																																																																																						        }
																																																																																																																																																																																																																																																																																								        tags = append(tags, tagName)
																																																																																																																																																																																																																																																																																										    }

																																																																																																																																																																																																																																																																																											    w.Header().Set("Content-Type", "application/json")
																																																																																																																																																																																																																																																																																												    json.NewEncoder(w).Encode(tags)
																																																																																																																																																																																																																																																																																													    log.Println("Fetched repository details successfully")
																																																																																																																																																																																																																																																																																														}

																																																																																																																																																																																																																																																																																														// PullImage pulls a Docker image from Docker Hub
																																																																																																																																																																																																																																																																																														func PullImage(w http.ResponseWriter, r *http.Request) {
																																																																																																																																																																																																																																																																																														    image := r.URL.Query().Get("image")
																																																																																																																																																																																																																																																																																															    log.Printf("Pulling image: %s\n", image)

																																																																																																																																																																																																																																																																																																    cmd := exec.Command("docker", "pull", image)
																																																																																																																																																																																																																																																																																																	    if len(settings.Accelerators) > 0 {
																																																																																																																																																																																																																																																																																																		        cmd.Env = append(cmd.Env, fmt.Sprintf("DOCKER_OPTS=--registry-mirror=%s", settings.Accelerators[0]))
																																																																																																																																																																																																																																																																																																				    }

																																																																																																																																																																																																																																																																																																					    if err := cmd.Run(); err != nil {
																																																																																																																																																																																																																																																																																																						        log.Printf("Error pulling image: %v\n", err)
																																																																																																																																																																																																																																																																																																								        http.Error(w, err.Error(), http.StatusInternalServerError)
																																																																																																																																																																																																																																																																																																										        return
																																																																																																																																																																																																																																																																																																												    }

																																																																																																																																																																																																																																																																																																													    w.Header().Set("Content-Type", "application/json")
																																																																																																																																																																																																																																																																																																														    json.NewEncoder(w).Encode(map[string]bool{"success": true})
																																																																																																																																																																																																																																																																																																															    log.Println("Image pulled successfully")
																																																																																																																																																																																																																																																																																																																}

																																																																																																																																																																																																																																																																																																																// UpdateSettings updates the accelerators settings
																																																																																																																																																																																																																																																																																																																func UpdateSettings(w http.ResponseWriter, r *http.Request) {
																																																																																																																																																																																																																																																																																																																    var newSettings Settings
																																																																																																																																																																																																																																																																																																																	    if err := json.NewDecoder(r.Body).Decode(&newSettings); err != nil {
																																																																																																																																																																																																																																																																																																																		        log.Printf("Error decoding settings: %v\n", err)
																																																																																																																																																																																																																																																																																																																				        http.Error(w, err.Error(), http.StatusBadRequest)
																																																																																																																																																																																																																																																																																																																						        return
																																																																																																																																																																																																																																																																																																																								    }

																																																																																																																																																																																																																																																																																																																									    settings = newSettings

																																																																																																																																																																																																																																																																																																																										    w.Header().Set("Content-Type", "application/json")
																																																																																																																																																																																																																																																																																																																											    json.NewEncoder(w).Encode(map[string]bool{"success": true})
																																																																																																																																																																																																																																																																																																																												    log.Println("Settings updated successfully")
																																																																																																																																																																																																																																																																																																																													}

																																																																																																																																																																																																																																																																																																																													// GetSettings gets the current accelerators settings
																																																																																																																																																																																																																																																																																																																													func GetSettings(w http.ResponseWriter, r *http.Request) {
																																																																																																																																																																																																																																																																																																																													    w.Header().Set("Content-Type", "application/json")
																																																																																																																																																																																																																																																																																																																														    json.NewEncoder(w).Encode(settings)
																																																																																																																																																																																																																																																																																																																															    log.Println("Fetched settings successfully")
																																																																																																																																																																																																																																																																																																																																}