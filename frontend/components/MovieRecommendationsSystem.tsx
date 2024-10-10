"use client";

import React, { useState, useEffect } from "react";
import { Search, Bookmark, BookmarkCheck } from "lucide-react";

interface MovieSuggestion {
  id: number;
  title_x: string;
}

interface Cast {
  credit_id: string;
  department: string;
  gender: number;
  id: number;
  job: string;
  name: string;
}

interface Genre {
  id: number;
  name: string;
}

interface Keyword {
  id: number;
  name: string;
}

interface ProductionCompany {
  id: number;
  name: string;
}

interface ProductionCountry {
  iso_3166_1: string;
  name: string;
}

interface SpokenLanguage {
  iso_639_1: string;
  name: string;
}

interface MovieResponse {
  _id: string;
  backdrop_path: string;
  budget: number;
  cast: Cast[];
  genres: Genre[];
  homepage: string;
  id: number;
  keywords: Keyword[];
  original_language: string;
  original_title: string;
  overview: string;
  popularity: number;
  production_companies: ProductionCompany[];
  production_countries: ProductionCountry[];
  release_date: string;
  revenue: number;
  runtime: number;
  spoken_languages: SpokenLanguage[];
  status: string;
  tagline: string;
  title_x: string;
  title_y: string;
  vote_average: number;
  vote_count: number;
}

interface Recommendation {
  adult: boolean;
  backdrop_path: string;
  genre_ids: number[];
  id: number;
  media_type: string;
  original_language: string;
  original_title: string;
  overview: string;
  popularity: number;
  poster_path: string;
  release_date: string;
  title: string;
  video: boolean;
  vote_average: number;
  vote_count: number;
}

interface RecommendationsResponse {
  page: number;
  results: Recommendation[];
  total_pages: number;
  total_results: number;
}

const MovieRecommendationSystem: React.FC = () => {
  const [searchTerm, setSearchTerm] = useState("");
  const [suggestions, setSuggestions] = useState<MovieSuggestion[]>([]);
  const [selectedMovie, setSelectedMovie] = useState<MovieResponse | null>(null);
  const [recommendations, setRecommendations] = useState<Recommendation[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [bookmarkedMovies, setBookmarkedMovies] = useState<Set<number>>(new Set());

  useEffect(() => {
    if (searchTerm.length > 2 && isSearching) {
      const delayDebounceFn = setTimeout(() => {
        fetchSuggestions(searchTerm);
      }, 300);

      return () => clearTimeout(delayDebounceFn);
    } else {
      setSuggestions([]);
    }
  }, [searchTerm, isSearching]);

  const fetchSuggestions = async (term: string) => {
    try {
      const response = await fetch(`http://localhost:8080/api/v1/essearch?value=${encodeURIComponent(term)}`);
      const data = await response.json();
      setSuggestions(data);
    } catch (error) {
      console.error("Error fetching suggestions:", error);
    }
  };

  const fetchMovieDetails = async (movieId: number) => {
    try {
      const response = await fetch(`http://localhost:8080/api/v1/searchmovie?movieId=${movieId}`);
      const data: MovieResponse = await response.json();
      setSelectedMovie(data);
    } catch (error) {
      console.error("Error fetching movie details:", error);
    }
  };

  const fetchRecommendations = async (movieId: number) => {
    try {
      const response = await fetch(`http://localhost:8080/api/v1/recommendations?movieId=${movieId}`);
      const data: RecommendationsResponse = await response.json();
      setRecommendations(data.results);
    } catch (error) {
      console.error("Error fetching recommendations:", error);
    }
  };

  const handleMovieSelection = async (movie: MovieSuggestion) => {
    setSearchTerm(movie.title_x);
    setSuggestions([]);
    setIsSearching(false);
    await fetchMovieDetails(movie.id);
    await fetchRecommendations(movie.id);
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(e.target.value);
    setIsSearching(true);
  };

  const toggleBookmark = (movieId: number) => {
    setBookmarkedMovies((prevBookmarks) => {
      const newBookmarks = new Set(prevBookmarks);
      if (newBookmarks.has(movieId)) {
        newBookmarks.delete(movieId);
      } else {
        newBookmarks.add(movieId);
      }
      return newBookmarks;
    });
  };

  const isBookmarked = (movieId: number) => bookmarkedMovies.has(movieId);

  return (
    <div className="container mx-auto p-4 text-black">
      <h1 className="text-2xl font-bold mb-4">Movie Recommendation System</h1>
      <div className="relative">
        <input
          type="text"
          value={searchTerm}
          onChange={handleInputChange}
          onFocus={() => setIsSearching(true)}
          className="w-full p-2 border rounded-md"
          placeholder="Search for a movie..."
        />
        <span className="absolute right-2 top-2 text-gray-500">
          <Search size={20} />
        </span>
        {suggestions && suggestions.length > 0 && isSearching && (
          <ul className="absolute z-10 w-full bg-white border rounded-md mt-1 max-h-60 overflow-auto">
            {suggestions.map((movie) => (
              <li
                key={movie.id}
                onClick={() => handleMovieSelection({ id: movie.id, title_x: movie.title_x })}
                className="p-2 hover:bg-gray-100 cursor-pointer"
              >
                {movie.title_x}
              </li>
            ))}
          </ul>
        )}
      </div>

      {selectedMovie && (
        <div className="mt-8 flex flex-col md:flex-row">
          <div className="md:w-1/3 mb-4 md:mb-0">
            <img
              src={`https://image.tmdb.org/t/p/w500${selectedMovie.backdrop_path}`}
              alt={`${selectedMovie.title_x} poster`}
              className="w-full rounded-md shadow-lg"
            />
          </div>
          <div className="md:w-2/3 md:pl-8">
            <div className="flex justify-between items-center mb-2">
              <h2 className="text-xl font-bold">{selectedMovie.title_x}</h2>
              <button onClick={() => toggleBookmark(selectedMovie.id)} className="text-blue-500 hover:text-blue-700">
                {isBookmarked(selectedMovie.id) ? <BookmarkCheck size={24} /> : <Bookmark size={24} />}
              </button>
            </div>
            <p className="mb-2">{selectedMovie.overview}</p>
            <p className="mb-2">Release Date: {selectedMovie.release_date}</p>
            <p className="mb-2">Rating: {selectedMovie.vote_average}/10</p>
            <h3 className="text-lg font-semibold mt-4 mb-2">Cast</h3>
            <ul className="list-disc pl-5">
              {selectedMovie.cast && selectedMovie.cast.slice(0, 5).map((actor) => <li key={actor.credit_id}>{actor.name}</li>)}
            </ul>
          </div>
        </div>
      )}

      {recommendations.length > 0 && (
        <div className="mt-8">
          <h2 className="text-xl font-bold mb-4">Recommended Movies</h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            {recommendations.map((movie) => (
              <div key={movie.id} className="border rounded-md p-4">
                <img
                  src={`https://image.tmdb.org/t/p/w500${movie.backdrop_path}`}
                  alt={`${movie.title} poster`}
                  className="w-full h-40 object-cover rounded-md mb-2"
                />
                <h3 className="font-semibold mb-2">{movie.title}</h3>
                <p className="text-sm">{movie.overview.slice(0, 100)}...</p>
                <button
                  onClick={() => handleMovieSelection({ id: movie.id, title_x: movie.title })}
                  className="mt-2 bg-blue-500 text-white px-2 py-1 rounded-md text-sm"
                >
                  View Details
                </button>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default MovieRecommendationSystem;
