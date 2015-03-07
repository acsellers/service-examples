class SuggestController < ApplicationController
  before_filter :validate_token


  def index
    posts = Post.where(published: true)

    if params["search"].present?
      posts = posts.where("name LIKE ?", "%#{search}%")
    end

    if params["recent"].present?
      posts = posts.where("created_at > ?", Time.now - 2.weeks)
    end

    if params["all"].blank?
      posts = posts.limit(15)
    end

    render json: posts
  end

  def create
    imported = 0
    params["post"]["posts"].each do |param_post|
      post = Post.find_by(permalink: param_post["permalink"]) || Post.new
      post.attributes = param_post
      imported+=1 if post.save
    end
    render json: imported
  end

  def delete
    Post.where(permalink: params[:id]).each(&:delete)
    render json: "ok"
  end

  def show
    exists = Post.where(permalink: params[:id]).exists?
    render json: { exists: exists }
  end

  private

  def validate_token
    token = request.headers["AUTH_TOKEN"]
    Rails.log.info("Token:", token)
    # check that the auth_token is valid
  end

end
