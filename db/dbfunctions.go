package db

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetPodcastByURL(url string, podcast *Podcast) error {
	result := DB.Preload(clause.Associations).Where(&Podcast{URL: url}).First(&podcast)
	return result.Error
}

func GetPodcastsByURLList(urls []string, podcasts *[]Podcast) error {
	result := DB.Preload(clause.Associations).Where("url in ?", urls).First(&podcasts)
	return result.Error
}
func GetAllPodcasts(podcasts *[]Podcast) error {

	result := DB.Preload("PodcastItems").Find(&podcasts)
	return result.Error
}
func GetAllPodcastItems(podcasts *[]PodcastItem) error {

	result := DB.Preload("Podcast").Order("pub_date desc").Find(&podcasts)
	return result.Error
}
func GetPaginatedPodcastItems(page int, count int, podcasts *[]PodcastItem, total *int64) error {
	result := DB.Debug().Preload("Podcast").Limit(count).Offset((page - 1) * count).Order("pub_date desc").Find(&podcasts)

	DB.Model(&PodcastItem{}).Count(total)

	return result.Error
}
func GetPodcastById(id string, podcast *Podcast) error {

	result := DB.Preload("PodcastItems", func(db *gorm.DB) *gorm.DB {
		return db.Order("podcast_items.pub_date DESC")
	}).First(&podcast, "id=?", id)
	return result.Error
}

func GetPodcastItemById(id string, podcastItem *PodcastItem) error {

	result := DB.Preload(clause.Associations).First(&podcastItem, "id=?", id)
	return result.Error
}
func DeletePodcastItemById(id string) error {

	result := DB.Where("id=?", id).Delete(&PodcastItem{})
	return result.Error
}
func DeletePodcastById(id string) error {

	result := DB.Where("id=?", id).Delete(&Podcast{})
	return result.Error
}

func GetAllPodcastItemsByPodcastId(podcastId string, podcasts *[]PodcastItem) error {

	result := DB.Preload(clause.Associations).Where(&PodcastItem{PodcastID: podcastId}).Find(&podcasts)
	return result.Error
}

func GetAllPodcastItemsToBeDownloaded() (*[]PodcastItem, error) {
	var podcastItems []PodcastItem
	result := DB.Debug().Preload(clause.Associations).Where("download_status=?", NotDownloaded).Find(&podcastItems)
	//fmt.Println("To be downloaded : " + string(len(podcastItems)))
	return &podcastItems, result.Error
}
func GetAllPodcastItemsAlreadyDownloaded() (*[]PodcastItem, error) {
	var podcastItems []PodcastItem
	result := DB.Debug().Preload(clause.Associations).Where("download_status=?", Downloaded).Find(&podcastItems)
	return &podcastItems, result.Error
}

func GetPodcastItemByPodcastIdAndGUID(podcastId string, guid string, podcastItem *PodcastItem) error {

	result := DB.Preload(clause.Associations).Where(&PodcastItem{PodcastID: podcastId, GUID: guid}).First(&podcastItem)
	return result.Error
}
func GetPodcastByTitleAndAuthor(title string, author string, podcast *Podcast) error {

	result := DB.Preload(clause.Associations).Where(&Podcast{Title: title, Author: author}).First(&podcast)
	return result.Error
}

func CreatePodcast(podcast *Podcast) error {
	tx := DB.Create(&podcast)
	return tx.Error
}

func CreatePodcastItem(podcastItem *PodcastItem) error {
	tx := DB.Omit("Podcast").Create(&podcastItem)
	return tx.Error
}
func UpdatePodcastItem(podcastItem *PodcastItem) error {
	tx := DB.Omit("Podcast").Save(&podcastItem)
	return tx.Error
}
func UpdateSettings(setting *Setting) error {
	tx := DB.Save(&setting)
	return tx.Error
}
func GetOrCreateSetting() *Setting {
	var setting Setting
	result := DB.First(&setting)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		DB.Save(&Setting{})
		DB.First(&setting)
	}
	return &setting
}
