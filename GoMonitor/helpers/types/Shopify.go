package Types

type ShopifyNewProduct struct {
	Store     string `json:"store,omitempty"`
	BodyHTML  string `json:"body_html,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Handle    string `json:"handle,omitempty"`
	ID        int64  `json:"id,omitempty"`
	Images    []struct {
		CreatedAt  string        `json:"created_at,omitempty"`
		Height     int64         `json:"height,omitempty"`
		ID         int64         `json:"id,omitempty"`
		Position   int64         `json:"position,omitempty"`
		ProductID  int64         `json:"product_id,omitempty"`
		Src        string        `json:"src,omitempty"`
		UpdatedAt  string        `json:"updated_at,omitempty"`
		VariantIds []interface{} `json:"variant_ids,omitempty"`
		Width      int64         `json:"width,omitempty"`
	} `json:"images,omitempty"`
	Options []struct {
		Name     string   `json:"name,omitempty"`
		Position int64    `json:"position,omitempty"`
		Values   []string `json:"values,omitempty"`
	} `json:"options,omitempty"`
	ProductType string    `json:"product_type,omitempty"`
	PublishedAt string    `json:"published_at,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	Title       string    `json:"title,omitempty"`
	UpdatedAt   string    `json:"updated_at,omitempty"`
	Variants    []Variant `json:"variants,omitempty"`
	Vendor      string    `json:"vendor,omitempty"`
}

type Variant struct {
	Available        bool        `json:"available,omitempty"`
	CompareAtPrice   string      `json:"compare_at_price,omitempty"`
	CreatedAt        string      `json:"created_at,omitempty"`
	FeaturedImage    interface{} `json:"featured_image,omitempty"`
	Grams            int64       `json:"grams,omitempty"`
	ID               int64       `json:"id,omitempty"`
	Option1          string      `json:"option1,omitempty"`
	Option2          string      `json:"option2,omitempty"`
	Option3          interface{} `json:"option3,omitempty"`
	Position         int64       `json:"position,omitempty"`
	Price            string      `json:"price,omitempty"`
	ProductID        int64       `json:"product_id,omitempty"`
	RequiresShipping bool        `json:"requires_shipping,omitempty"`
	Sku              string      `json:"sku,omitempty"`
	Taxable          bool        `json:"taxable,omitempty"`
	Title            string      `json:"title,omitempty"`
	UpdatedAt        string      `json:"updated_at,omitempty"`
}
type ShopifyProductJS struct {
	Available            bool     `json:"available,omitempty"`
	Store                string   `json:"store,omitempty"`
	CompareAtPrice       int64    `json:"compare_at_price,omitempty"`
	CompareAtPriceMax    int64    `json:"compare_at_price_max,omitempty"`
	CompareAtPriceMin    int64    `json:"compare_at_price_min,omitempty"`
	CompareAtPriceVaries bool     `json:"compare_at_price_varies,omitempty"`
	CreatedAt            string   `json:"created_at,omitempty"`
	Description          string   `json:"description,omitempty"`
	FeaturedImage        string   `json:"featured_image,omitempty"`
	Handle               string   `json:"handle,omitempty"`
	ID                   int64    `json:"id,omitempty"`
	Images               []string `json:"images,omitempty"`
	Media                []struct {
		Alt          interface{} `json:"alt,omitempty"`
		AspectRatio  float64     `json:"aspect_ratio,omitempty"`
		Height       int64       `json:"height,omitempty"`
		ID           int64       `json:"id,omitempty"`
		MediaType    string      `json:"media_type,omitempty"`
		Position     int64       `json:"position,omitempty"`
		PreviewImage struct {
			AspectRatio float64 `json:"aspect_ratio,omitempty"`
			Height      int64   `json:"height,omitempty"`
			Src         string  `json:"src,omitempty"`
			Width       int64   `json:"width,omitempty"`
		} `json:"preview_image,omitempty"`
		Src   string `json:"src,omitempty"`
		Width int64  `json:"width,omitempty"`
	} `json:"media,omitempty"`
	Options []struct {
		Name     string   `json:"name,omitempty"`
		Position int64    `json:"position,omitempty"`
		Values   []string `json:"values,omitempty"`
	} `json:"options,omitempty"`
	Price       int64    `json:"price,omitempty"`
	PriceMax    int64    `json:"price_max,omitempty"`
	PriceMin    int64    `json:"price_min,omitempty"`
	PriceVaries bool     `json:"price_varies,omitempty"`
	PublishedAt string   `json:"published_at,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Title       string   `json:"title,omitempty"`
	Type        string   `json:"type,omitempty"`
	URL         string   `json:"url,omitempty"`
	Variants    []struct {
		Available           bool        `json:"available,omitempty"`
		Barcode             string      `json:"barcode,omitempty"`
		CompareAtPrice      int64       `json:"compare_at_price,omitempty"`
		FeaturedImage       interface{} `json:"featured_image,omitempty"`
		ID                  int64       `json:"id,omitempty"`
		InventoryManagement string      `json:"inventory_management,omitempty"`
		InventoryPolicy     string      `json:"inventory_policy,omitempty"`
		InventoryQuantity   int64       `json:"inventory_quantity,omitempty"`
		Name                string      `json:"name,omitempty"`
		Option1             string      `json:"option1,omitempty"`
		Option2             string      `json:"option2,omitempty"`
		Option3             interface{} `json:"option3,omitempty"`
		Options             []string    `json:"options,omitempty"`
		Price               int64       `json:"price,omitempty"`
		PublicTitle         string      `json:"public_title,omitempty"`
		RequiresShipping    bool        `json:"requires_shipping,omitempty"`
		Sku                 string      `json:"sku,omitempty"`
		Taxable             bool        `json:"taxable,omitempty"`
		Title               string      `json:"title,omitempty"`
		Weight              int64       `json:"weight,omitempty"`
	} `json:"variants,omitempty"`
	Vendor string `json:"vendor,omitempty"`
}

type SiteCollectMongo struct {
	Store    string
	Products []ShopifyNewProduct
}

type ShopifyJson struct {
	Products []ShopifyNewProduct
}
type ProductsItem struct {
	ID          int64
	UpdatedTime string
	Variants    []Variant
	Store       string
	Handle      string
}
