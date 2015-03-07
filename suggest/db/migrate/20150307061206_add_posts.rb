class AddPosts < ActiveRecord::Migration
  def change
    create_table :posts do |t|
      t.string :name
      t.text :body
      t.string :permalink
      t.timestamps
      t.bool :published
    end
  end
end
